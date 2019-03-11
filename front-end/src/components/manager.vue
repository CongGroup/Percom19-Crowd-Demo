<template>
    <div>
        <div class="agg">
            <div class="container" v-if="!enableStatics">
                <div class="admin">
                    <div class="contract">
                        <div class="role">
                            <div class ="item">
                                <div class="stage"> Data Consumer </div>
                            </div>
                        </div>
                        <div class="item">
                            <div class="label">Account: </div>
                            <div class="value"> {{dataConsumerAccount.address}}</div>
                        </div>
                        <div class="item">
                            <div class="label">Balance: </div>
                            <div class="value"> {{consumerTokenBalance}} </div>
                        </div>
                        <div class="col buttonGroup">
                            <div class="form" v-if="atStage('solicit')">
                                <div v-if="solicitStatus === 0">
                                    <span class="label">Target Number:</span>
                                    <input class="input" v-model.number="targetNumber">
                                    <button  :disabled="!atStage('solicit')" class="btn btn-dark contract-button" @click="solicit"> Solicit</button>
                                </div>
                                <div class="buttonWaiting" v-else-if="solicitStatus===1">
                                    <div class="label">
                                         Soliciting
                                    </div>
                                    <pacman></pacman>
                                </div>
                            </div>
                            <div class="form" v-if="atStage('register')">
                                <button :disabled="!atStage('register')" v-if="stopRegisterStatus === 0" class="btn btn-dark contract-button" @click="stopRegisterAndSubmit"> Stop Registration</button>
                                <div class="buttonWaiting" v-else-if="stopRegisterStatus===1">
                                    <div class="label">
                                        Stopping
                                    </div>
                                    <pacman></pacman>
                                </div>
                            </div>
                            <div class="row">
                                <div class="form" v-if="atStage('approve')">
                                    <button :disabled="!atStage('approve')" v-if="approveStatus === 0" class="btn btn-dark contract-button" @click="approve"> Approve</button>
                                    <div class="buttonWaiting" v-else-if="approveStatus === 1">
                                        <div class="label">
                                            Approving
                                            <pacman></pacman>
                                        </div>
                                    </div>
                                </div>
                                <div class="form" v-if="atStage('claim')">
                                    <button :disabled="false" v-if="!showingStatus" class="btn btn-dark contract-button" @click="showStatics"> Show Statics</button>
                                    <pacman v-else></pacman>
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
                            <div class="label">Account: </div>
                            <div class="value">{{serviceProviderAccount.address}}  </div>
                        </div>
                        <div class="item">
                            <div class="label">Balance: </div>
                            <div class="value"> {{tokenBalance}} </div>
                        </div>
                        <div class="col buttonGroup">
                            <div class="form" v-if="atStage('aggregate')">
                                <button :disabled="!atStage('aggregate')" v-if="aggregateStatus === 0" class="btn btn-dark contract-button" @click="aggregate"> Aggregate</button>
                                <div class="buttonWaiting" v-else-if="aggregateStatus === 1">
                                    <div class="label">
                                        Aggregating
                                    </div>
                                    <pacman></pacman>
                                </div>
                            </div>
                            <div class="form" v-if="atStage('claim')">
                                <button :disabled="!atStage('claim')" v-if="claimStatus === 0" class="btn btn-dark contract-button" @click="claim"> Claim</button>
                                <div class="buttonWaiting" v-else-if="claimStatus === 1">
                                    <div class="label">
                                        Claiming
                                    </div>
                                    <pacman></pacman>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="adminStage">
                        <div class="stageForm">
                            <img class="icon" alt="Vue logo" src="../assets/cd3.jpeg"/>
                            <!--<div class="value">{{stageToProcedure[stage]}} </div>-->
                        </div>
                        <div class="stageHead">
                            Task Information
                        </div>
                        <div class="stageContent contract">
                            <div >
                                <div class="item">
                                    <div class="label">Reward for data provider: </div>
                                    <div class="value"> {{solicitInfo.dataFee}} </div>
                                </div>
                                <div class="item">
                                    <div class="label">Reward for service provider: </div>
                                    <div class="value"> {{solicitInfo.serviceFee}} </div>
                                </div>
                                <div class="item">
                                    <div class="label">Expected number of submissions: </div>
                                    <div class="value"> {{solicitInfo.target}}</div>
                                </div>
                            </div>
                            <div>
                                <div class="item">
                                    <div class="label">Valid data range </div>
                                    <div class="value" v-if="shouldShow('register')"> 0—65535</div>
                                </div>
                            </div>
                            <div class="item" >
                                <span class="label">Current number of submissions: </span>
                                <span class="value"> {{registerNumber}} </span>
                            </div>
                            <!--<div class="item" >-->
                            <!--<span class="label">Submission number:</span>-->
                            <!--<span class="value"> {{submissionNumber}} </span>-->
                            <!--</div>-->
                            <div class="item" >
                                <span class="label">Number of valid submissions::</span>
                                <span class="value"> {{qualifiedNumber}} </span>
                            </div>
                            <!--<div class="item" >-->
                            <!--<span class="label">Final aggregate result:</span>-->
                            <!--<span class="value"> {{qualifiedNumber !==0 ?aggregateResult:"NAN"}} </span>-->
                            <!--</div>-->
                            <div class="item">
                                <span class="label">Claim number:</span>
                                <span class="value"> {{claimNumber}} </span>
                            </div>
                        </div>
                </div>
            </div>
            <div class="graphLayout"  v-else >
                <div class="head">
                    Task Statistics
                </div>
                <div class="contract info">
                    <!--<div class="item" >-->
                    <!--<span class="label">Register number: </span>-->
                    <!--<span class="value"> {{registerNumber}} </span>-->
                    <!--</div>-->
                    <!--<div class="item" >-->
                    <!--<span class="label">Submission number:</span>-->
                    <!--<span class="value"> {{submissionNumber}} </span>-->
                    <!--</div>-->
                    <div class="item" >
                        <span class="label" >Number of valid submissions:</span>
                        <span class="value"> {{qualifiedNumber}} </span>
                    </div>
                    <div class="item" >
                        <span class="label" >Aggregation sum:</span>
                        <span class="value"> {{qualifiedNumber !==0 ?aggregateResult:"NAN"}} </span>
                    </div>
                    <div class="item" >
                        <span class="label" >Aggregation average:</span>
                        <span class="value"> {{qualifiedNumber !==0 ?Math.floor(aggregateResult/qualifiedNumber):"NAN"}} </span>
                    </div>
                    <!--<div class="item">-->
                    <!--<span class="label">Claim number:</span>-->
                    <!--<span class="value"> {{claimNumber}} </span>-->
                    <!--</div>-->
                </div>
                <div class="plotHeader">
                    Distribution of valid submissions
                </div>
                <div class="graph">
                    <!--<canvas  id="myChart"></canvas>-->
                    <bar-statistics :width="900" :height="600" :chart-data="graph" :options="graphOptions"></bar-statistics>
                </div>
                <div class="labels">
                    <div class="label" v-for="(l,i) in labels">
                        <span v-if="i<=1">{{l}} </span>
                        <span v-else-if="i===2" style="margin-left: -15px">{{l}}</span>
                        <span style="margin-left: -16px;" v-else-if="i>=2">{{l}}</span>
                    </div>
                </div>
                <div class="samplesHeader">
                    Invalid submission samples
                </div>
                <div class="samples">
                    <div class="sample" v-for="(s,i) in invalidSamples">
                        <div class="value">
                            {{i}}.
                        </div>
                        <div>
                            <span style="font-weight: 600;">{{s}}</span>
                        </div>
                    </div>
                </div>
                <div class="back">
                    <button class="btn btn-dark contract-button" @click="hidenGraph"> Back </button>
                </div>
            </div>
        </div>
    </div>
</template>


<script>
    import barChart from './bar.js';
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

    // const GCUID_CAN_REGISTER_AND_SUBMIT_FOR_TASK = 800;
    const GCUID_CAN_CLAIM_FOR_TASK = 801;

    const TASK_ID = 0;

    const INITIAL = 0;
    const WAITING = 1;
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
                labels: undefined,
                reconnecting: false,
                invalidSamples: undefined,
                showingStatus: false,
                solicitStatus: 2,
                stopRegisterStatus:2,
                aggregateStatus:2,
                approveStatus:2,
                claimStatus:2,
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
                wsReconnect:undefined,
            }
        },
        computed: {
        },
        watch: {
        },
        methods: {
            hidenGraph: function(event){
                this.enableStatics = false;
            },
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
                this.solicitStatus = WAITING;
                console.log("solicit");
                let payload = {
                    gcuid: GCUID_SOLICIT,
                    dataFee: 10,
                    serviceFee: 50,
                    serviceProvider: this.serviceProviderAccount.address,
                    target: this.targetNumber,
                    privateKey: this.dataConsumerAccount.privateKey,
                    address: this.dataConsumerAccount.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            aggregate: function() {
                this.aggregateStatus = WAITING;
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
                this.approveStatus = WAITING;
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
                this.claimStatus = WAITING;
                console.log("claim");
                let payload = {
                    gcuid: GCUID_CLAIM,
                    taskId: TASK_ID,
                    privateKey: this.serviceProviderAccount.privateKey,
                    address: this.serviceProviderAccount.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            canClaim: function() {
                console.log("query can claim");
                let payload = {
                    gcuid: GCUID_CAN_CLAIM_FOR_TASK,
                    taskId: TASK_ID,
                    address: this.serviceProviderAccount.address,
                };
                this.ws.send(JSON.stringify(payload));
            },
            stopRegisterAndSubmit: function () {
                this.stopRegisterStatus = WAITING;
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
            draw: function() {
                let qualifiedData = this.submitValues;
                console.log("qualified data", this.submitValues)
                let space = 8192;
                let minH = 0;
                let maxH = 65536;
                let bucket = Array((maxH-minH)/space).fill(0);
                qualifiedData.forEach(v=>{
                    let bucketNumber = Math.floor((v-minH)/space);
                    // if(bucketNumber >= bucket.length) {
                    //     bucketNumber = bucket.length - 1;
                    // }
                    // if(bucketNumber < 0) {
                    //     bucketNumber = 0;
                    // }
                    ++bucket[bucketNumber];
                });

                bucket.forEach((v,i)=>{
                    console.log(`${i},${v}`);
                    bucket[i] = (v/qualifiedData.length).toFixed(2);
                });

                let labels = ['0'];
                let labels2 = [''];
                for(let i=0;i<bucket.length;++i) {
                    labels.push(`${space*(i+1)}`);
                    labels2.push('');
                }
                this.labels = labels;
                labels = labels2;

                let datasets = [
                    {
                        label: 'Height',
                        data: bucket,
                        borderColor: '#002266',
                        backgroundColor:'#002266',
                        hoverBackgroundColor:'#ff5050',
                        hoverBorderColor: '#ff5050',
                        hoverBorderWidth: 1,
                    },
                ];

                this.graphOptions = {
                    legend: {
                        display: false,
                    },
                    scales: {
                        yAxes: [{
                            ticks: {
                                beginAtZero: true,
                                max:1.0,
                                fontSize: 15,
                            },
                        }],
                        xAxes: [{
                            // gridLines: {
                            //     offset: true
                            // },
                            type:'category',
                            labels: labels,
                            ticks: {
                                fontSize: 15,
                            },
                        }],
                    },
                };
                this.graph = {
                    // labels: labels,
                    datasets: datasets,
                };
            },
            showStatics: function() {
                this.showingStatus = true;
                console.log("showStatics");
                if(this.graph===undefined) {
                    if(this.submitValues===undefined) {
                        console.log("refectch submit values");
                        let p = this.axios.get(`${process.env.HTTP_PATH}/statistics/${TASK_ID}`);
                        p.then((res)=>{
                            this.submitValues = res.data.submitValues;
                            this.invalidSamples = res.data.invalidSamples;
                            console.log("submit values", this.submitValues);
                            console.log("invalid samples:",this.invalidSamples);
                            this.draw();
                            this.enableStatics = true;
                            this.showingStatus = false;
                        })
                    } else {
                        this.draw();
                        this.enableStatics = true;
                        this.showingStatus = false;
                    }
                } else {
                    this.enableStatics = true;
                    this.showingStatus = false;
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
                this.invalidSamples = undefined;
                this.enableStatics = false;
                this.solicitInfo = {
                    dataFee: undefined,
                    serviceFee: undefined,
                    serviceProvider: undefined,
                    target: undefined
                }
            },
            initialDisplay: function() {
                this.requireEther();
                this.requireDataConsumerEther();
                this.getCurrentStage();
                this.getTokenBalance();
                this.getConsumerTokenBalance();
            },
            requireEther: function() {
                this.axios.get(`${process.env.HTTP_PATH}/requireEther/${this.serviceProviderAccount.address}`).then(res=>{
                    let balance = res.data;
                    if(balance>=ETHER_THREASHOLD){
                        this.hasEther = true;
                        console.log("already has ether");
                    } else {
                        console.log("do not have enough ether, get ether now");
                        this.axios.get(`${process.env.HTTP_PATH}/requireEther/${this.serviceProviderAccount.address}`).then(res=>{
                            let value = res.data;
                            console.log(`get ${value} ether`);
                            this.hasEther = true;
                        }).catch(err=>{
                            console.log(err.message);
                        })
                    }
                })
            },
            requireDataConsumerEther: function() {
                this.axios.get(`${process.env.HTTP_PATH}/requireEther/${this.dataConsumerAccount.address}`).then(res=>{
                    let balance = res.data;
                    if(balance>=ETHER_THREASHOLD){
                        this.hasEther = true;
                        console.log("already has ether");
                    } else {
                        console.log("do not have enough ether, get ether now");
                        this.axios.get(`${process.env.HTTP_PATH}/requireEther/${this.dataConsumerAccount.address}`).then(res=>{
                            let value = res.data;
                            console.log(`get ${value} ether`);
                            this.hasEther = true;
                        }).catch(err=>{
                            console.log(err.message);
                        })
                    }
                })
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
                this.ws = new WebSocket(process.env.SERVER_PATH);
                this.ws.onopen = e => {
                    console.log("websocket open");
                    this.reconnecting = false;
                    this.initialDisplay();
                };
                this.ws.onmessage = e => {
                    let res = JSON.parse(e.data);
                    switch (res.gcuid) {
                        case GCUID_SOLICIT:
                            if (res.status === 0) {
                                this.getConsumerTokenBalance();
                                this.solicitStatus = SUBMITTED;
                            } else {
                                console.log(res.reason);
                                this.solicitStatus = INITIAL;
                                // TODO ERROR WARNING
                            }
                            break;
                        case GCUID_AGGREGATION:
                            if (res.status === 0) {
                                this.aggregateStatus = SUBMITTED;
                            } else {
                                console.log(res.reason);
                                this.aggregateStatus = INITIAL;
                                // TODO ERROR WARNING
                            }
                            break;
                        case GCUID_APPROVE:
                            if (res.status === 0) {
                                // this.submitValues = res.submitValues;
                                // this.invalidSamples =res.invalidSamples;
                                // console.log("invalid samples:",this.invalidSamples);
                                // console.log("submit values:",this.submitValues);
                                this.approveStatus = SUBMITTED;
                            } else {
                                console.log(res.reason);
                                this.approveStatus = INITIAL;
                                // TODO ERROR WARNING
                            }
                            break;
                        case GCUID_CLAIM:
                            if (res.status === 0) {
                                this.claimStatus = SUBMITTED;
                                this.getTokenBalance();
                            } else {
                                this.claimStatus = INITIAL;
                                console.log(res.reason);
                                // TODO ERROR WARNING
                            }
                            break;
                        case GCUID_STOP_REGISTER_AND_SUBMIT:
                            if (res.status === 0) {
                                this.stopRegisterStatus = SUBMITTED;
                            } else {
                                this.stopRegisterStatus = INITIAL;
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
                                if(this.submissionNumber === undefined || this.submissionNumber<res.amount)
                                this.submissionNumber = res.amount;
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_REGISTER_NUMBER:
                            if (res.status === 0) {
                                if(this.registerNumber===undefined || this.registerNumber<res.amount)
                                this.registerNumber = res.amount;
                            } else {
                                console.log(res.reason);
                            }
                            break;
                        case GCUID_CURRENT_STAGE:
                            if (res.status === 0) {
                                this.stage = res.stage;
                                // console.log(this.stageToProcedure[this.stage]);
                                if(this.stageToProcedure[this.stage] === 'Solicit') {
                                    this.solicitStatus = INITIAL;
                                    this.cleanState();
                                } else {
                                    if(!this.initialied) {
                                        console.log("get info");
                                        this.getInfo();
                                        this.initialied = true;
                                    } else {
                                        this.initializeState(res.stage);
                                    }
                                    if (this.stageToProcedure[this.stage] === 'Register') {
                                        this.stopRegisterStatus = INITIAL;
                                    } else if (this.stageToProcedure[this.stage] === 'Aggregate') {
                                        this.aggregateStatus = INITIAL;
                                    }  else if(this.stageToProcedure[this.stage] === 'Approve') {
                                        this.approveStatus = INITIAL;
                                    } else if(this.stageToProcedure[this.stage] === 'Claim') {
                                        this.canClaim()
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
                                    if(this.claimNumber === undefined || this.claimNumber<res.amount) this.claimNumber = res.amount;
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
                        case GCUID_CAN_CLAIM_FOR_TASK:
                            console.log("can calim?:",res.canClaim);
                            if (res.status === 0) {
                                if(res.canClaim) {
                                    this.claimStatus = INITIAL;
                                } else {
                                    this.claimStatus = SUBMITTED;
                                }
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
                    // TODO ERROR WARNING
                    this.reconnecting = true;
                    this.reconnect()
                };
            },
            reconnect: function() {
               this.wsReconnect = setTimeout(this.initialWS,2000);
            },
        },
        created: function () {
            this.initialWS();
            console.log(this.ws)
        },
        beforeDestroy: function() {
            if(this.wsReconnect!==undefined) clearTimeout(this.wsReconnect)
        }
    };
</script>
<!-- Add "scoped" attribute to limit CSS to this component only -->

<style scoped>
    /*@import '../assets/css/style.css';*/
</style>
