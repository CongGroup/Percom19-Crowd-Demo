<template>
    <div class="agg">
        <div class="display">
            <div class="contract">
                <img class="icon" alt="Vue logo" src="../assets/cd3.jpeg"/>
                <div class="item">
                    <div class="label">Your Address: </div>
                    <div class="value" v-if="account!==undefined">{{address}}  </div>
                </div>
                <div class="item">
                    <div class="label">Your Balance: </div>
                    <div class="value"> {{ether}} </div>
                </div>
                <div class="item">
                    <span class="label">Current Stage:  </span>
                    <span class="value">{{stageToProcedure[stage]}} </span>
                </div>
                <div class="item">
                    <span class="label">Register number: </span>
                    <span class="value"> {{registerNumber}} </span>
                </div>
                <div class="item">
                    <span class="label">Submission number:</span>
                    <span class="value"> {{submissionNumber}} </span>
                </div>
                <div class="item">
                    <span class="label">Final aggregate result:</span>
                    <span class="value"> </span>
                </div>
                <div class="item">
                    <span class="label">Claim number:</span>
                    <span class="value"> {{claimNumber}} </span>
                </div>
            </div>
        </div>

        <div v-if="account!==undefined">
            <div class="form" v-if="atStage('solicit')">
                <span class="label">Value:</span>
                <input class="input" v-model.number="value">
                <button class="btn btn-dark" @click="solicit"> solicit</button>
            </div>
            <div v-else-if="equalsToSubmitState(1)">
                <div class="waiting">Soliciting {{waiting}}</div>
            </div>
        </div>

        <div v-if="account!==undefined">
            <div class="form" v-if="atStage('register')">
                <span class="label">Value:</span>
                <input class="input" v-model.number="value">
                <button class="btn btn-dark" @click="register"> register</button>
            </div>
            <div v-else-if="equalsToSubmitState(1)">
                <div class="waiting">Submitting {{waiting}}</div>
            </div>
        </div>

        <div v-if="account!==undefined">
            <div class="form" v-if="atStage('submit')">
                <span class="label">Value:</span>
                <input class="input" v-model.number="value">
                <button class="btn btn-dark" @click="submit"> submit</button>
            </div>
            <div v-else-if="equalsToSubmitState(1)">
                <div class="waiting">Submitting {{waiting}}</div>
            </div>
        </div>

        <div v-if="account!==undefined">
            <div class="form" v-if="atStage('aggregate')">
                <span class="label">Value:</span>
                <input class="input" v-model.number="value">
                <button class="btn btn-dark" @click="aggregate"> aggregate</button>
            </div>
            <div v-else-if="equalsToSubmitState(1)">
                <div class="waiting">Submitting {{waiting}}</div>
            </div>
        </div>

        <div v-if="account!==undefined">
            <div class="form" v-if="atStage('approve')">
                <span class="label">Value:</span>
                <input class="input" v-model.number="value">
                <button class="btn btn-dark" @click="approve"> approve</button>
            </div>
            <div v-else-if="equalsToSubmitState(1)">
                <div class="waiting">Submitting {{waiting}}</div>
            </div>
        </div>

        <div v-if="account!==undefined">
            <div class="form" v-if="atStage('claim')">
                <span class="label">Value:</span>
                <input class="input" v-model.number="value">
                <button class="btn btn-dark" @click="claim"> claim</button>
            </div>
            <div v-else-if="equalsToSubmitState(1)">
                <div class="waiting">Submitting {{waiting}}</div>
            </div>
        </div>

        <div class="waiting" v-else>
            Initializing wallet {{waiting}}
        </div>
    </div>
</template>


<script>
    const GCUID_SOLICIT = 0;
    const GCUID_REGISTER = 1;
    const GCUID_SUBMIT = 2;
    const GCUID_AGGREGATION = 3;
    const GCUID_APPROVE = 4;
    const GCUID_CLAIM = 5;

    const GCUID_ETHER = 101;
    const GCUID_CURRENT_STAGE = 102;
    const GCUID_REGISTER_NUMBER = 103;
    const GCUID_SUBMISSION_NUMBER = 104;
    const GCUID_SOLICIT_INFO = 105;
    const GCUID_AGGREGATE_RESULT = 106;
    const GCUID_CLAIM_NUMBER = 107;

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
                submitStatus: INITIAL,
                waiting: "",
                waitingAnimate: undefined,
                ws: undefined,
                value: 0,
                ether: undefined,
                stage: undefined,
                registerNumber: undefined,
                submissionNumber: undefined,
                claimNumber: undefined,
                mapToStage : {
                    "solicit":0,
                    "register":1,
                    "submit":2,
                    "aggregate":3,
                    "approve":4,
                    "claim": 5,
                },
                stageToProcedure : {
                    0: "solicit",
                    1: "register",
                    2: "submit",
                    3: "aggregate",
                    4: "approve",
                    5: "claim",
                },
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
                    this.endWaiting();
                    this.getEther();
                }
            }
        },
        methods: {
            atStage: function (stage) {
                // return true;
                  return this.stage === this.mapToStage[stage];
            },
            getEther: function () {
                console.log("get Balance");
                let payload = {
                    "gcuid": GCUID_ETHER,
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
                // TODO
            },
            getAggregateResult: function() {
                // TODO
            },
            solicit: function() {
                console.log("solicit");
                let payload = {
                    gcuid: GCUID_SOLICIT,
                    dataFee: 1000,
                    serviceFee: 10,
                    serviceProvider: this.account.address,
                    target: 1,
                    privateKey: this.account.privateKey,
                    address: this.account.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            register: function() {
                console.log("register");
                let payload = {
                    gcuid: GCUID_REGISTER,
                    taskId: TASK_ID,
                    privateKey: this.account.privateKey,
                    address: this.account.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            submit: function() {
                console.log("submitting");
                this.startWaiting();
                this.submitStatus = SUBMITTING;

                let payload = {
                    gcuid: GCUID_SUBMIT,
                    taskId: TASK_ID,
                    value: this.value,
                    address: this.account.address,
                    privateKey: this.account.privateKey,
                };

                this.ws.send(JSON.stringify(payload));

                this.submitStatus = SUBMITTED;
                this.endWaiting();
            },
            aggregate: function() {
                console.log("aggregate");
                let payload = {
                    gcuid: GCUID_AGGREGATION,
                    taskId: TASK_ID,
                    privateKey: this.account.privateKey,
                    address: this.account.address
                };
                this.ws.send(JSON.stringify(payload));
            },
            approve: function() {
                console.log("approve");
                let payload = {
                    gcuid: GCUID_APPROVE,
                    taskId: TASK_ID,
                    privateKey: this.account.privateKey,
                    address: this.account.address
                };
                this.ws.send(JSON.stringify(payload));
            },
            claim: function() {
                console.log("claim");
                let payload = {
                    gcuid: GCUID_CLAIM,
                    taskId: TASK_ID,
                    privateKey: this.account.privateKey,
                    address: this.account.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            equalsToSubmitState: function (state) {
                return this.submitStatus === state ;
            },
            startWaiting: function () {
                this.waitingAnimate = setInterval(()=>{
                    if(this.waiting.length>=4) this.waiting = "";
                    else this.waiting += '.';
                },1000);
            },
            endWaiting: function () {
                clearInterval(this.waitingAnimate);
                this.waiting="";
            },
            formatEther: function(ether) {
                return parseFloat(ether)/10**18
            },
            initialDisplay: function() {
                this.getRegisterNumber();
                this.getSubmissionNumber();
                this.getCurrentStage();
                this.getClaimNumber();
            }
        },
        created: function () {
            this.startWaiting();
            this.ws = new WebSocket("ws://localhost:4000");
            this.ws.onopen = e => {
                console.log("websocket open");
                this.initialDisplay();
            };
            this.ws.onmessage = e => {
                let res = JSON.parse(e.data);
                console.log(res);
                switch (res.gcuid) {
                    case GCUID_SOLICIT:
                        if (res.status === 0) {
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
                        } else {
                            console.log(res.reason);
                        }
                        break;
                    case GCUID_CLAIM:
                        if (res.status === 0) {
                        } else {
                            console.log(res.reason);
                        }
                        break;
                    case GCUID_ETHER:
                        if (res.status === 0) {
                            this.ether = this.formatEther(res.amount);
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
                        } else {
                            console.log(res.reason);
                        }
                        break;
                    case GCUID_SOLICIT_INFO:
                        // TODO
                        if (res.status === 0) {
                            this.stage = res.stage;
                        } else {
                            console.log(res.reason);
                        }
                        break;
                    case GCUID_AGGREGATE_RESULT:
                        // TODO
                        if (res.status === 0) {
                            this.claimNumber = res.amount;
                        } else {
                            console.log(res.reason);
                        }
                        break;
                    case GCUID_CLAIM_NUMBER:
                        if (res.status === 0) {
                            this.stage = res.stage;
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
            };
        }
    };
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->

<style scoped>
    /*@import '../assets/css/style.css';*/
</style>
