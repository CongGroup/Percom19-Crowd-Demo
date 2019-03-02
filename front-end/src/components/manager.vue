<template>
    <div class="agg">
        <div class="container">
            <div class="adminStage">
                <div class="stageForm">
                    <img class="icon" alt="Vue logo" src="../assets/cd3.jpeg"/>
                    <div class="value">{{stageToProcedure[stage]}} </div>
                </div>
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
                        <span class="label">Final aggregate result:</span>
                        <span class="value"> {{aggregateResult}} </span>
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

                    <div class="col buttonGroup">
                        <div class="form" v-if="atStage('solicit')">
                            <span class="label">Target Number:</span>
                            <input class="input" v-model.number="value">
                            <button  :disabled="!atStage('solicit')" class="btn btn-dark contract-button" @click="solicit"> solicit</button>
                        </div>
                        <div class="form">
                            <button :disabled="!atStage('approve')" class="btn btn-dark contract-button" @click="approve"> approve</button>
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
                        <div class="value"> {{ether}} </div>
                    </div>
                    <div class="col buttonGroup">
                        <div class="form">
                            <button :disabled="!atStage('aggregate')" class="btn btn-dark contract-button" @click="aggregate"> aggregate</button>
                        </div>
                        <div class="form">
                            <button :disabled="!atStage('claim')" class="btn btn-dark contract-button" @click="claim"> claim</button>
                        </div>
                    </div>
                </div>
            </div>
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
        name: "manager",
        props: {
            serviceProviderAccount: Object,
            dataConsumerAccount: Object,
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
                initialied: false,
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
                return this.stage == this.mapToStage[stage];
            },
            getEther: function () {
                console.log("get Balance:",this.serviceProviderAccount.address);
                let payload = {
                    "gcuid": GCUID_ETHER,
                    "address": this.serviceProviderAccount.address
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
                    target: 1,
                    privateKey: this.dataConsumerAccount.privateKey,
                    address: this.dataConsumerAccount.address,
                };
                this.ws.send(JSON.stringify(payload));
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
            formatEther: function(ether) {
                return parseFloat(ether)/10**18
            },
            cleanState: function() {
                this.claimNumber = undefined;
                this.registerNumber = undefined;
                this.submissionNumber = undefined;
                this.aggregateResult = undefined;
                this.solicitInfo = {
                    dataFee: undefined,
                    serviceFee: undefined,
                    serviceProvider: undefined,
                    target: undefined
                }
            },
            initialDisplay: function() {
                this.getCurrentStage();
                this.getEther()
            },
            initializeState: function(stage) {
                if('register' === this.mapToStage[stage]) {
                    this.registerNumber = 0;
                } else if('submit' === this.mapToStage[stage]) {
                    this.submissionNumber = 0;
                } else if('claim' === this.mapToStage[stage]) {
                    this.claimNumber = 0;
                }
            },
            getInfo: function() {
                if(this.shouldShow('approve')) {
                    this.getAggregateResult();
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
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_CLAIM_NUMBER:
                            if (res.status === 0) {
                                this.claimNumber = res.amount;
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
        created: function () {
            this.initialWS();
            console.log(this.ws)
        }
    };
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->

<style scoped>
    /*@import '../assets/css/style.css';*/
</style>
