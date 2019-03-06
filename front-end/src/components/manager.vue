<template>
    <div class="agg">
        <div class="container">
            <div class="adminStage">
                <div class="stageForm" v-if="!enableStatics">
                    <img class="icon" alt="Vue logo" src="../assets/cd3.jpeg"/>
                    <div class="value">{{stageToProcedure[stage]}} </div>
                </div>
                <bar-statistics :chart-data="graph" :options="graphOptions" v-else></bar-statistics>
                <div class="stageContent contract">
                    <div >
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
                    <div class="item" >
                        <span class="label">Register number: </span>
                        <span class="value"> {{registerNumber}} </span>
                    </div>
                    <div class="item" >
                        <span class="label">Submission number:</span>
                        <span class="value"> {{submissionNumber}} </span>
                    </div>
                    <div class="item" >
                        <span class="label">Qualified number:</span>
                        <span class="value"> {{qualifiedNumber}} </span>
                    </div>
                    <div class="item" >
                        <span class="label">Final aggregate result:</span>
                        <span class="value"> {{qualifiedNumber !==0 ?aggregateResult:"NAN"}} </span>
                    </div>
                    <div class="item">
                        <span class="label">Claim number:</span>
                        <span class="value"> {{claimNumber}} </span>
                    </div>
                </div>
            </div>
            <div class="admin">
                <div class="contract">
                    <div class="role">
                        <div class ="item">
                            <div class="stage"> Data Consumer </div>
                        </div>
                    </div>
                    <div class="item">
                        <div class="label">Address: </div>
                        <div class="value"> {{dataConsumerAccount.address}}</div>
                    </div>
                    <div class="item">
                        <div class="label">Your Balance: </div>
                        <div class="value"> {{consumerTokenBalance}} </div>
                    </div>
                    <div class="col buttonGroup">
                        <div class="form" v-if="atStage('solicit')">
                            <span class="label">Target Number:</span>
                            <input class="input" v-model.number="targetNumber">
                            <button  :disabled="!atStage('solicit')" class="btn btn-dark contract-button" @click="solicit"> solicit</button>
                        </div>
                        <div class="form" v-if="atStage('register')">
                            <button :disabled="!atStage('register')" class="btn btn-dark contract-button" @click="stopRegisterAndSubmit"> stop register</button>
                        </div>
                        <div class="row">
                            <div class="form" v-if="atStage('approve')">
                                <button :disabled="!atStage('approve')" class="btn btn-dark contract-button" @click="approve"> approve</button>
                            </div>
                            <div class="form" v-if="atStage('claim')">
                                <button :disabled="false" class="btn btn-dark contract-button" @click="showStatics"> showStatics</button>
                            </div>
                        </div>
                    </div>

                </div>
                <div class="contract">
                    <div class="role">
                        <div class ="item">
                            <div class="stage">Service Provider </div>
                        </div>
                    </div>
                    <div class="item">
                        <div class="label">Address: </div>
                        <div class="value">{{serviceProviderAccount.address}}  </div>
                    </div>
                    <div class="item">
                        <div class="label">Your Balance: </div>
                        <div class="value"> {{tokenBalance}} </div>
                    </div>
                    <div class="col buttonGroup">
                        <div class="form" v-if="atStage('aggregate')">
                            <button :disabled="!atStage('aggregate')" class="btn btn-dark contract-button" @click="aggregate"> aggregate</button>
                        </div>
                        <div class="form" v-if="atStage('claim')">
                            <button :disabled="!atStage('claim')" class="btn btn-dark contract-button" @click="claim"> claim</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>


<script>
    import barChart from './bar.js';
    import {getBaseTxObject,encodeFunction,signTx} from '../assets/js/tx.js'
    const GCUID_SOLICIT = 0;
    const GCUID_REGISTER = 1;
    const GCUID_SUBMIT = 2;
    const GCUID_AGGREGATION = 3;
    const GCUID_APPROVE = 4;
    const GCUID_CLAIM = 5;

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
    const GCUID_SUBMIT_VALUES = 111;

    const TASK_ID = 0;

    const INITIAL = 0;
    const SUBMITTING = 1;
    const SUBMITTED = 2;

    export default {
        name: "manager",
        props: {
            serviceProviderAccount: Object,
            dataConsumerAccount: Object,
        },
        components: {
            barStatistics: barChart
        },
        data: function () {
            return {
                wsPath: "ws://0.0.0.0:4000",
                httpPath: "http://0.0.0.0:4000",
                ws: undefined,
                value: 0,
                tokenBalance: undefined,
                consumerTokenBalance: undefined,
                stage: undefined,
                registerNumber: undefined,
                submissionNumber: undefined,
                submitValues: undefined,
                claimNumber: undefined,
                initialied: false,
                qualifiedNumber: undefined,
                enableStatics: false,
                graph: undefined,
                graphOptions: undefined,
                mapToStage : {
                    "solicit":0,
                    "register":1,
                    "submit":2,
                    "aggregate":3,
                    "approve":4,
                    "claim": 5,
                },
                targetNumber: 1,
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
            }
        },
        computed: {
        },
        watch: {
        },
        methods: {
            shouldShow: function (stage) {
                return this.stage >= this.mapToStage[stage];
            },
            atStage: function(stage) {
                return this.stage === this.mapToStage[stage];
            },
            getTokenBalance: function () {
                console.log("get Balance:",this.serviceProviderAccount.address);
                let payload = {
                    "gcuid": GCUID_BALANCE,
                    "address": this.serviceProviderAccount.address
                };
                this.ws.send(JSON.stringify(payload));
            },
            getConsumerTokenBalance: function() {
                let payload = {
                    "gcuid": GCUID_BALANCE,
                    "address": this.dataConsumerAccount.address
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
            solicit: function() {
                console.log("solicit");
                let payload = {
                    gcuid: GCUID_SOLICIT,
                    dataFee: 1230,
                    serviceFee: 16,
                    serviceProvider: this.serviceProviderAccount.address,
                    target: this.targetNumber,
                    privateKey: this.dataConsumerAccount.privateKey,
                    address: this.dataConsumerAccount.address,
                };
                this.ws.send(JSON.stringify(payload));
                // send by mobile
                // let data = encodeFunction("solicit",1230,16,this.serviceProviderAccount.address,this.targetNumber);
                // let p1 = this.axios.get(`${this.httpPath}/nonce/${this.dataConsumerAccount.address}`);
                // let p2 = this.axios.get(`${this.httpPath}/chainId`);
                // Promise.all([p1,p2]).then(([r1,r2])=>{
                //     let nonce = r1.data;
                //     let chainId = r2.data;
                //     console.log(`nonce:${nonce}`);
                //     console.log(`chainId:${chainId}`);
                //     let tx = getBaseTxObject();
                //     tx.nonce = nonce;
                //     tx.data = data;
                //     console.log(`data:${data}`);
                //     tx.chainId = chainId;
                //     return signTx(tx,'0x'+this.dataConsumerAccount.privateKey)
                // }).then((rawTx)=>{
                //     let payload = {
                //         gcuid: GCUID_SEND_TRANSACTION,
                //         rawTransaction: rawTx.slice(2),
                //     };
                //     this.ws.send(JSON.stringify(payload));
                // }).catch(err=>console.log(err))
            },
            aggregate: function() {
                console.log("aggregate");
                let payload = {
                    gcuid: GCUID_AGGREGATION,
                    taskId: TASK_ID,
                    privateKey: this.serviceProviderAccount.privateKey,
                    address: this.serviceProviderAccount.address
                };
                this.ws.send(JSON.stringify(payload));
            },
            approve: function() {
                console.log("approve");
                let payload = {
                    gcuid: GCUID_APPROVE,
                    taskId: TASK_ID,
                    privateKey: this.dataConsumerAccount.privateKey,
                    address: this.dataConsumerAccount.address
                };
                this.ws.send(JSON.stringify(payload));
            },
            claim: function() {
                console.log("claim");
                let payload = {
                    gcuid: GCUID_CLAIM,
                    taskId: TASK_ID,
                    privateKey: this.serviceProviderAccount.privateKey,
                    address: this.serviceProviderAccount.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            stopRegisterAndSubmit: function () {
                console.log("stop register and submit");
                let payload = {
                    gcuid: GCUID_STOP_REGISTER_AND_SUBMIT,
                    taskId:  TASK_ID,
                    privateKey: this.dataConsumerAccount.privateKey,
                    address: this.dataConsumerAccount.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            getQualifiedNumber: function() {
                console.log("get qualifiedNumber");
                let payload = {
                    gcuid: GCUID_QUALIFIED_NUMBER,
                    taskId: TASK_ID,
                };
                this.ws.send(JSON.stringify(payload));
            },
            showStatics: function() {
                console.log("showStatics");

                if(this.graph===undefined) {
                    if(this.submitValues===undefined) {
                        console.log("refectch submit values");
                        let p = this.axios.get(`${this.httpPath}/statistics/${TASK_ID}`);
                        p.then((res)=>{
                            this.submitValues = res.data.submitValues;
                            console.log("submitValues:",this.submitValues)
                        })
                    }

                    let qualifiedData = [154,165,170,156,182,183,190,166,165,160,124,200,212,323];
                    let space = 5;
                    let minH = 150;
                    let maxH = 200;
                    let bucket = Array((maxH-minH)/space+2).fill(0);
                    qualifiedData.forEach(v=>{
                        let bucketNumber = Math.floor((v-minH)/space)+1;
                        if(bucketNumber >= bucket.length) {
                            bucketNumber = bucket.length - 1;
                        }
                        if(bucketNumber < 0) {
                            bucketNumber = 0;
                        }
                        ++bucket[bucketNumber];
                    });

                    let labels = [];
                    for(let i=0;i<bucket.length;++i) {
                        if(i===0) labels.push(`0-${minH}`);
                        else if(i===bucket.length-1) labels.push(`>${maxH}`);
                        else {
                            let start = minH + space*(i-1);
                            let end = minH + space*i;
                            labels.push(`${start}-${end}`);
                        }
                    }

                    let datasets = [
                        {
                            label: 'Height',
                            data: bucket,
                            borderColor: '#002266',
                            backgroundColor:'#002266',
                            hoverBackgroundColor:'#ff5050',
                            hoverBorderColor: '#ff5050',
                            hoverBorderWidth: 1,
                        },{
                            label:'',
                            data: bucket,
                            type:'line',
                            backgroundColor:'#ffffff',
                            borderColor: '#000000',
                        }
                    ];

                    this.graph = {
                        labels: labels,
                        datasets: datasets,
                    };
                }
                if(!this.enableStatics) {
                    this.enableStatics = true;
                } else {
                    this.enableStatics = false;
                }
            },
            cleanState: function() {
                console.log("clean state!!!!!!");
                this.claimNumber = undefined;
                this.registerNumber = undefined;
                this.submissionNumber = undefined;
                this.aggregateResult = undefined;
                this.qualifiedNumber = undefined;
                this.submitValues = undefined;
                this.enableStatics = false;
                this.solicitInfo = {
                    dataFee: undefined,
                    serviceFee: undefined,
                    serviceProvider: undefined,
                    target: undefined
                }
            },
            initialDisplay: function() {
                this.getCurrentStage();
                this.getTokenBalance();
                this.getConsumerTokenBalance();
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
                this.ws = new WebSocket(this.wsPath);
                this.ws.onopen = e => {
                    console.log("websocket open");
                    this.initialDisplay();
                };
                this.ws.onmessage = e => {
                    let res = JSON.parse(e.data);
                    switch (res.gcuid) {
                        case GCUID_SOLICIT:
                            if (res.status === 0) {
                                this.getConsumerTokenBalance();
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_REGISTER:
                            if (res.status === 0) {
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_SUBMIT:
                            if (res.status === 0) {
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_AGGREGATION:
                            if (res.status === 0) {
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_APPROVE:
                            if (res.status === 0) {
                                this.submitValues = res.submitValues;
                                console.log("submit values:",this.submitValues);
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
                        case GCUID_STOP_REGISTER_AND_SUBMIT:
                            if (res.status === 0) {
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_BALANCE:
                            if (res.status === 0) {
                                if(res.address === this.serviceProviderAccount.address) {
                                    this.tokenBalance =  res.amount;
                                } else if(res.address === this.dataConsumerAccount.address) {
                                    this.consumerTokenBalance = res.amount;
                                }
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
                                        this.initializeState(res.stage);
                                    }
                                }
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_SEND_TRANSACTION:
                            if (res.status === 0) {
                                console.log("send transaction")
                            } else {
                                console.log(res.reason)
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
                                console.log("aggregate result",res)
                                this.aggregateResult = res.amount;
                                this.qualifiedNumber = res.qualifiedNumber;
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_CLAIM_NUMBER:
                            if (res.status === 0) {
                                console.log("claim number:"+res.amount);
                                if(this.stage === this.mapToStage['claim']) {
                                    this.claimNumber = res.amount;
                                }
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_QUALIFIED_NUMBER:
                            if (res.status === 0) {
                                this.qualifiedNumber = res.amount;
                                console.log("qualified number:",this.qualifiedNumber)
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
