<template>
    <div class="agg">
        <div class="display">
            <div class="contract">
                <div class ="item">
                    <img class="icon" alt="Vue logo" src="../assets/cd3.jpeg"/>
                    <div class="stage">{{stageToProcedure[stage]}} stage </div>
                </div>
                <div class="item">
                    <div class="label">Your Address: </div>
                    <div class="value" v-if="account!==undefined">{{address}}  </div>
                </div>
                <div class="item">
                    <div class="label">Your Balance: </div>
                    <div class="value"> {{tokenBalance}} </div>
                </div>
                <div v-if="shouldShow('register')">
                    <div class="item">
                        <div class="label">Solicit Data Fee: </div>
                        <div class="value"> {{solicitInfo.dataFee}} </div>
                    </div>
                    <div class="item">
                        <div class="label">Solicit Service Fee: </div>
                        <div class="value"> {{solicitInfo.serviceFee}} </div>
                    </div>
                    <div class="item">
                        <div class="label">Solicit Target Number: </div>
                        <div class="value"> {{solicitInfo.target}}</div>
                    </div>
                </div>
                <div class="item" v-if="shouldShow('register')">
                    <span class="label">Register number: </span>
                    <span class="value"> {{registerNumber}} </span>
                </div>
                <div class="item" v-if="shouldShow('submit')">
                    <span class="label">Submission number:</span>
                    <span class="value"> {{submissionNumber}} </span>
                </div>
                <div class="item" v-if="shouldShow('approve')">
                    <span class="label">Qualified Number:  </span>
                    <span class="value">{{qualifiedNumber}}</span>
                </div>
                <div class="item" v-if="shouldShow('approve')">
                    <span class="label">Are you qualified?  </span>
                    <span class="value">{{qualified === undefined ? '': qualified?'true':'false'}}</span>
                </div>
                <div class="item" v-if="shouldShow('approve')">
                    <span class="label">Final aggregate result:</span>
                    <span class="value"> {{ qualifiedNumber !==0 ?aggregateResult:"NAN"}} </span>
                </div>
                <div class="item" v-if="shouldShow('claim')">
                    <span class="label">Claim number:</span>
                    <span class="value"> {{claimNumber}} </span>
                </div>
            </div>
        </div>

        <!--<div v-if="account!==undefined">-->
            <!--<div class="form" v-if="atStage('register')">-->
                <!--<span class="label">Value:</span>-->
                <!--<input class="input" v-model.number="value">-->
                <!--<button class="btn btn-dark" @click="register"> register</button>-->
            <!--</div>-->
        <!--</div>-->

        <div v-if="account!==undefined">
            <div class="form" v-if="atStage('register')">
                <span class="label">Value:</span>
                <input class="input" v-model.number="value">
                <button class="btn btn-dark" @click="registerAndSubmit"> submit</button>
            </div>
        </div>


        <div class="formLists" v-if="account!==undefined">
            <div class="form" v-if="atStage('claim')">
                <button class="btn btn-dark" @click="claim"> claim</button>
            </div>
        </div>

        <div class="waiting" v-else>
            Initializing wallet {{waiting}}
        </div>
    </div>
</template>


<script>
    import {getBaseTxObject,signTx,encodeFunction} from "../assets/js/tx";

    const GCUID_SOLICIT = 0;
    const GCUID_REGISTER = 1;
    const GCUID_SUBMIT = 2;
    const GCUID_AGGREGATION = 3;
    const GCUID_APPROVE = 4;
    const GCUID_CLAIM = 5;
    const GCUID_REGISTER_AND_SUBMIT = 6;
    const GCUID_STOP_REGISTER_AND_SUBMIT = 7;
    const GCUID_SEND_TRANSACTION = 8;

    const GCUID_ETHER = 101;
    const GCUID_CURRENT_STAGE = 102;
    const GCUID_REGISTER_NUMBER = 103;
    const GCUID_SUBMISSION_NUMBER = 104;
    const GCUID_SOLICIT_INFO = 105;
    const GCUID_AGGREGATE_RESULT = 106;
    const GCUID_CLAIM_NUMBER = 107;
    const GCUID_BALANCE = 108;
    const GCUID_QUALIFIED_NUMBER = 109;
    const GCUID_IS_QUALIFIED = 110;

    const TASK_ID = 0;

    const INITIAL = 0;
    const SUBMITTING = 1;
    const SUBMITTED = 2;

    export default {
        name: "HelloWorld",
        props: {
            account: Object
        },
        data: function () {
            return {
                wsPath: "ws://0.0.0.0:4000",
                httpPath: "http://0.0.0.0:4000",
                initialized:false,
                submitStatus: INITIAL,
                waiting: "",
                waitingAnimate: undefined,
                ws: undefined,
                value: 0,
                tokenBalance: undefined,
                stage: undefined,
                registerNumber: undefined,
                submissionNumber: undefined,
                claimNumber: undefined,
                qualified: undefined,
                mapToStage : {
                    "solicit":0,
                    "register":1,
                    "submit":2,
                    "aggregate":3,
                    "approve":4,
                    "claim": 5,
                },
                stageToProcedure : {
                    0: "Solicit",
                    1: "Register",
                    2: "Submit",
                    3: "Aggregate",
                    4: "Approve",
                    5: "Claim",
                },
                solicitInfo : {
                    dataFee: undefined,
                    serviceFee: undefined,
                    serviceProvider: undefined,
                    target: undefined
                },
                aggregateResult: undefined,
                qualifiedNumber: undefined,
            }
        },
        computed: {
            address: function () {
                let tmp = this.account.address;
                return tmp.slice(0,8)+'...'+ tmp.slice(-6);
            }
        },
        watch: {
            account: function(newV,oldV){
                if(newV !== undefined){
                    this.getTokenBalance();
                }
            }
        },
        methods: {
            shouldShow: function (stage) {
                  return this.stage >= this.mapToStage[stage];
            },
            atStage: function(stage) {
                return this.stage == this.mapToStage[stage];
            },
            getTokenBalance: function () {
                console.log("get Balance");
                let payload = {
                    "gcuid": GCUID_BALANCE,
                    "address": this.account.address
                };
                this.ws.send(JSON.stringify(payload));
            },
            getCurrentStage: function() {
                console.log("get current stage");
                let payload = {
                    "gcuid": GCUID_CURRENT_STAGE,
                    "taskId": TASK_ID,
                };
                this.ws.send(JSON.stringify(payload));
            },
            getSubmissionNumber: function() {
                console.log("get current submission number");
                let payload = {
                    "gcuid": GCUID_SUBMISSION_NUMBER,
                    "taskId": TASK_ID,
                };
                this.ws.send(JSON.stringify(payload));
            },
            getRegisterNumber: function() {
                console.log("get current register number");
                let payload = {
                    "gcuid": GCUID_REGISTER_NUMBER,
                    "taskId": TASK_ID,
                };
                this.ws.send(JSON.stringify(payload));
            },
            getClaimNumber: function() {
                console.log("get current claim number");
                let payload = {
                    "gcuid": GCUID_CLAIM_NUMBER,
                    "taskId": TASK_ID,
                };
                this.ws.send(JSON.stringify(payload));
            },
            getSolicitInfo: function() {
                console.log("get solicit info");
                let payload = {
                    "gcuid": GCUID_SOLICIT_INFO,
                    "taskId": TASK_ID,
                };
                this.ws.send(JSON.stringify(payload));
            },
            getAggregateResult: function() {
                console.log("get aggregate result");
                let payload = {
                    "gcuid": GCUID_AGGREGATE_RESULT,
                    "taskId": TASK_ID,
                };
                this.ws.send(JSON.stringify(payload));
            },
            registerAndSubmit: function() {
                console.log("register and submit");
                // let payload = {
                //     gcuid: GCUID_REGISTER_AND_SUBMIT,
                //     taskId: TASK_ID,
                //     value: this.value,
                //     address: this.account.address,
                //     privateKey: this.account.privateKey,
                // };
                // this.ws.send(JSON.stringify(payload));
                console.log("address:"+this.account.address);
                let p1 = this.axios.get(`${this.httpPath}/nonce/${this.account.address}`);
                let p2 = this.axios.get(`${this.httpPath}/chainId`);
                let p3 = this.axios.post(`${this.httpPath}/encryptedData`,this.value);
                Promise.all([p1,p2,p3]).then(([r1,r2,r3])=>{
                    let nonce = r1.data;
                    let chainId = r2.data;
                    let encryptedData = r3.data.submitData;
                    let proof = r3.data.submitProof;
                    console.log(`nonce:${nonce}`);
                    console.log(`chainId:${chainId}`);
                    // console.log(`submitPayload:${JSON.stringify(r3.data)}`);

                    let data = encodeFunction("registerAndSubmit",TASK_ID,encryptedData,proof);
                    let tx = getBaseTxObject();
                    tx.from = this.account.address;
                    tx.nonce = nonce;
                    tx.data = data;
                    console.log(`data:${data}`);
                    tx.chainId = chainId;
                    return signTx(tx,'0x'+this.account.privateKey)
                }).then((rawTx)=>{
                    let payload = {
                        gcuid: GCUID_SEND_TRANSACTION,
                        rawTransaction: rawTx.slice(2),
                    };
                    this.ws.send(JSON.stringify(payload));
                }).catch(err=>console.log(err))
            },
            claim: function() {
                console.log("claim");
                // let payload = {
                //     gcuid: GCUID_CLAIM,
                //     taskId: TASK_ID,
                //     privateKey: this.account.privateKey,
                //     address: this.account.address,
                // };
                // this.ws.send(JSON.stringify(payload));
                let p1 = this.axios.get(`${this.httpPath}/nonce/${this.account.address}`);
                let p2 = this.axios.get(`${this.httpPath}/chainId`);
                Promise.all([p1,p2]).then(([r1,r2])=>{
                    let nonce = r1.data;
                    let chainId = r2.data;
                    console.log(`nonce:${nonce}`);
                    console.log(`chainId:${chainId}`);

                    let data = encodeFunction("claim",TASK_ID);
                    let tx = getBaseTxObject();
                    tx.nonce = nonce;
                    tx.data = data;
                    tx.from = this.accout;
                    console.log(`data:${data}`);
                    tx.chainId = chainId;
                    return signTx(tx,'0x'+this.account.privateKey)
                }).then((rawTx)=>{
                    let payload = {
                        gcuid: GCUID_SEND_TRANSACTION,
                        rawTransaction: rawTx.slice(2),
                    };
                    this.ws.send(JSON.stringify(payload));
                }).catch(err=>console.log(err))
            },
            getQualifiedNumber: function() {
                console.log("get qualifiedNumber");
                let payload = {
                    gcuid: GCUID_QUALIFIED_NUMBER,
                    taskId: TASK_ID,
                };
                this.ws.send(JSON.stringify(payload));
            },
            isQualified: function() {
                console.log("get isQualified");
                let payload = {
                    gcuid: GCUID_IS_QUALIFIED,
                    taskId: TASK_ID,
                    address: this.account.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            cleanState: function() {
                this.claimNumber = undefined;
                this.registerNumber = undefined;
                this.submissionNumber = undefined;
                this.qualifiedNumber = undefined;
                this.qualified = undefined;
                this.solicitInfo = {
                    dataFee: undefined,
                    serviceFee: undefined,
                    serviceProvider: undefined,
                    target: undefined
                }
            },
            initialDisplay: function() {
                this.getCurrentStage();
                if(this.account!==undefined) {
                    this.getTokenBalance();
                }
            },
            initializeState: function(stage) {
                if('register' === this.mapToStage[stage] && this.registerNumber === undefined) {
                    this.registerNumber = 0;
                } else if('submit' === this.mapToStage[stage] && this.submissionNumber=== undefined) {
                    this.submissionNumber = 0;
                } else if('claim' === this.mapToStage[stage] && this.claimNumber === undefined) {
                    this.claimNumber = 0;
                }
            },
            getInfo: function() {
                if(this.shouldShow('approve')) {
                    this.getAggregateResult();
                    this.getQualifiedNumber();
                    this.isQualified();
                }
                if(this.shouldShow('register')) {
                    this.getSolicitInfo();
                }
                if(this.shouldShow('claim')) {
                    this.getClaimNumber();
                }
                if(this.shouldShow('register')) {
                    this.getRegisterNumber();
                }
                if(this.shouldShow('submit')) {
                    this.getSubmissionNumber()
                }
            },
            initialWS: function() {
                this.ws = new WebSocket("ws://0.0.0.0:4000");
                this.ws.onopen = e => {
                    console.log("websocket open");
                    this.initialDisplay();
                };
                this.ws.onmessage = e => {
                    let res = JSON.parse(e.data);
                    switch (res.gcuid) {
                        case GCUID_REGISTER_AND_SUBMIT:
                            if (res.status === 0) {
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_CLAIM:
                            if (res.status === 0) {
                                this.getTokenBalance();
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_BALANCE:
                            if (res.status === 0) {
                                this.tokenBalance = res.amount;
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_SUBMISSION_NUMBER:
                            if (res.status === 0) {
                                this.submissionNumber = res.amount;
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_REGISTER_NUMBER:
                            if (res.status === 0) {
                                this.registerNumber = res.amount;
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_CURRENT_STAGE:
                            if (res.status === 0) {
                                this.stage = res.stage;
                                if(this.stageToProcedure[this.stage] === 'Solicit') {
                                    this.cleanState();
                                } else {
                                    if(!this.initialied) {
                                        console.log("get info");
                                        this.getInfo();
                                        this.initialied = true;
                                    } else {
                                        this.initializeState(res.stage)
                                    }
                                }
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_SOLICIT_INFO:
                            if (res.status === 0) {
                                this.solicitInfo = {
                                    dataFee: res.dataFee,
                                    serviceFee: res.serviceFee,
                                    serviceProvider: res.serviceProvider,
                                    target: res.target,
                                };
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_AGGREGATE_RESULT:
                            if (res.status === 0) {
                                this.aggregateResult = res.amount;
                                this.qualifiedNumber = res.qualifiedNumber;
                                this.isQualified();
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_CLAIM_NUMBER:
                            if (res.status === 0) {
                                console.log("claim number:",res.amount);
                                if(this.stage === this.mapToStage['claim']) {
                                    this.claimNumber = res.amount;
                                }
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_SEND_TRANSACTION:
                            if (res.status === 0) {
                                console.log("send transaction successfully");
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_QUALIFIED_NUMBER:
                            if (res.status === 0) {
                                this.qualifiedNumber = res.amount;
                                console.log("qualified number:",this.qualifiedNumber);
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_IS_QUALIFIED:
                            if (res.status === 0) {
                                this.qualified = res.qualified;
                                console.log("send transaction successfully");
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        default:
                            console.log("unknown response")
                    }
                };
                this.ws.onclose = e=> {
                    console.log("websocket close");
                    this.initialied = false;
                    this.reconnect()
                };
            },
            reconnect: function() {
                setTimeout(this.initialWS,2000);
            },
        },
        beforeMount: function () {
            this.initialWS();
            console.log(this.ws)
        },
    };
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->

<style scoped>
    /*@import '../assets/css/style.css';*/
</style>
