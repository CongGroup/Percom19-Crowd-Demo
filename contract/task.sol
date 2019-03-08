pragma solidity ^0.5.4;

import './IERC20.sol';

contract Crowdsourcing {

	uint constant public MAX_TASK = 1;
	uint constant internal TASK_FULL = 10;
	
	IERC20 public token;
	address private admin;

	enum Stages {
		solicit,
		register,
		submit,
		aggregate,
		approve,
		claim
	}

	struct Request {
		uint data_fee;
		uint service_fee;
		uint target;
	}


    struct DataProvider {
    	address account;
    	bool registered;
    	bool submited;
    	bool qualified;
    	bool claimed;
    	bytes submit_proof;
    	bytes submit_data;
    }

    struct ServiceProvider {
    	address account;
    	bool claimed;
    }

    struct Aggregation {
    	bytes aggregation;
    	bytes share;
    	bytes attestatino;
    }

	struct Task {
		Request  request;
		Stages  stage;
		address  owner;
		ServiceProvider service_provider;
		DataProvider [] data_provider;
		mapping(address => uint) data_provider_id;
		Aggregation aggregate;
		bool busy;
		uint register_count;
		uint submit_count;
		uint qualified_count;
		uint claim_count;
	}	
    
    uint public lastest_task;
    //Stages public lastest_stage = Stages.solicit;
	Task [MAX_TASK] internal task;

    
	// event Solicit(uint data_fee, uint service_fee, address owner, address service_provider, uint target, uint request_id, bytes32 task_id);
	// event Register(address data_provider, bytes32 task_id);
	// event RegisterCollected(bytes32 task_id);
	// event Submit(address data_provider,  bytes32 task_id);
	// event SubmitCollected(bytes32 task_id);
	// event Aggregate(bytes32 task_id);
	// event Approve(bytes32 task_id);
	// event Claim(address user, uint amount ,bytes32 task_id);
	// event ReceiveFund(uint amount, address supporter);
	event StageTransfer(uint  taskId, uint  newStage, uint  oldStage);
	event Solicit(uint  taskId);
	event Register(uint  taskId,uint registerNumber);
	event Submit(uint  taskId, uint submitNumber);
	event Aggregate(uint  taskId,bytes aggregateResult);
	event Approve(uint  taskId);
	event Claim(uint taskId,uint claimNumber);

	constructor (address erc20) public payable {
		token = IERC20(erc20);
		admin = msg.sender;
	}
	
	function changeERC20Contract(address erc20) external {
	    require(msg.sender==admin);
	    token = IERC20(erc20);
	}
	
	function () external payable {
	    //emit ReceiveFund(msg.value, msg.sender);
	}
	
	function getStageOfTask(uint task_id) external view returns(uint) {
	    return uint(task[task_id].stage);
	}
	
	function getSolicitInfoOfTask(uint task_id) external view returns(uint,uint,address,uint) {
	    return (task[task_id].request.data_fee,task[task_id].request.service_fee,task[task_id].service_provider.account,task[task_id].request.target);
	}
	
	function getSubmitDataOfTask(uint task_id, uint index) external view returns (bytes memory) {
	    return task[task_id].data_provider[index].submit_data;
	}
	
	function getSubmitCountOfTask(uint task_id) external view returns (uint) {
	    return task[task_id].submit_count;
	}
	
	function getRegisterNumberOfTask(uint task_id) external view returns (uint) {
	    return task[task_id].register_count;
	}
	
	function getClaimNumberOfTask(uint task_id) external view returns (uint) {
	    return task[task_id].claim_count;
	}
	
	function getAggregationResultOfTask(uint task_id) external view returns (bytes memory) {
	    return task[task_id].aggregate.aggregation;
	}
	

    function getSubmitProofOfTask(uint task_id, uint index) external view returns (bytes memory) {
        return task[task_id].data_provider[index].submit_proof;
    }
	

	function isQualifiedProviderForTask(uint task_id, address user) public view returns (bool) {
	    uint id = task[task_id].data_provider_id[user];
	    require(id!=0,"data provider not exist");
	    require(task[task_id].stage> Stages.aggregate,"only be called after aggregation!");
	    return task[task_id].data_provider[id-1].qualified;
	}
	
	function isRegisteredOfTask(uint task_id, address user) public view returns (bool) {
	    require(afterStage(task_id,Stages.solicit));
		uint id = task[task_id].data_provider_id[user];
		uint lastest_id = task[task_id].register_count;
		if(id == 0 || id > lastest_id || task[task_id].data_provider[id-1].account != user) {
		    return false;
		}
		return true;
	}
	
	function isSubmittedOfTask(uint task_id,address user) public view returns (bool) {
	    require(afterStage(task_id,Stages.register));
	    
	    if(!isRegisteredOfTask(task_id,user)) return false;
	    uint id = task[task_id].data_provider_id[user];
		if (!task[task_id].data_provider[id-1].submited) return false;
		
		return true;
	}
	
	function isServiceProviderOfTask(uint task_id, address user) public view returns (bool) {
	    if (user!=task[task_id].service_provider.account) return false;
	    return true;
	}
	
	function isClaimmedOfTask(uint task_id,address user) external view returns (bool) {
	    require(afterStage(task_id,Stages.approve));
	    if(!isServiceProviderOfTask(task_id,user) || !isSubmittedOfTask(task_id,user) || !isQualifiedProviderForTask(task_id,user)) return false;
	   
	    //service provider
        if(isServiceProviderOfTask(task_id,user)) {
            if(task[task_id].service_provider.claimed) return true;
            else return false;
        } 
        
        // data_provider
        uint id = task[task_id].data_provider_id[user];
        if(task[task_id].data_provider[id-1].claimed) return true;
	    return false;
	}
	

	function getQualifiedNumberOfTask(uint task_id) external view returns (uint) {
	    return task[task_id].qualified_count;
	}
	
	function setQualifiedDataProvider(uint task_id,bytes memory qualifiedSets) internal  {
	    uint qualified_count = 0;
	    for(uint i=0; i< qualifiedSets.length; ++i) {
	        if(qualifiedSets[i]!=bytes1(0)) {
	            qualified_count +=1;
	            task[task_id].data_provider[i].qualified = true;
	        } else {
	            task[task_id].data_provider[i].qualified = false;
	        }
	    }
	    task[task_id].qualified_count = qualified_count;
	}

	function getEmptyTaskSlot () internal view returns (uint) {
		for(uint i = 0; i< MAX_TASK ; ++i) {
			if (task[i].busy == false){
				return i;
			}
		}
		//not found
		return TASK_FULL;
	}

	function atStage (uint task_id, Stages _stage) internal view returns (bool) {
		if(task[task_id].stage == _stage){
			return true;
		} else {
			return false;
		}
	}
	
	function afterStage(uint task_id,Stages _stage) internal view returns (bool) {
		if(task[task_id].stage > _stage){
			return true;
		} else {
			return false;
		}
	}

	function nextStage(uint task_id) internal {
	    Stages oldStage = task[task_id].stage;
		if (task[task_id].stage == Stages.claim) {
			task[task_id].stage = Stages.solicit;
		} else {
	    	task[task_id].stage = Stages(uint(task[task_id].stage) + 1);
		}
		emit StageTransfer(task_id, uint(task[task_id].stage), uint(oldStage));
		//lastest_stage = task[task_id].stage;
	}

	function solicit(uint data_fee, uint service_fee, address service_provider, uint target) external {
		//require(address(this).balance > data_fee + service_fee);
		token.transferFrom(msg.sender,address(this),service_fee+data_fee);  // transfer token to contract 
		uint task_id = getEmptyTaskSlot();
		//lastest_task = task_id;
		require(task_id != TASK_FULL);
		task[task_id].request = Request(data_fee, service_fee, target);
		task[task_id].owner = msg.sender;
		task[task_id].service_provider = ServiceProvider(service_provider, false);
		task[task_id].busy = true;
		emit Solicit(task_id);
		nextStage(task_id);
		//emit Solicit(data_fee, service_fee, msg.sender, service_provider, target, request_id, task_id);
	}
    
	function register(uint task_id) external {
		require(atStage(task_id, Stages.register));
		address provider = msg.sender;
		uint id = task[task_id].data_provider_id[provider];
		uint lastest_id = task[task_id].register_count;
		require(id == 0 || id > lastest_id || task[task_id].data_provider[id-1].account != provider);
		if (task[task_id].data_provider.length == lastest_id) {
		    task[task_id].data_provider.push(DataProvider(provider,true,false,false,false,"0x1","0x1"));
		} else{
		    task[task_id].data_provider[lastest_id].account = provider;
		    task[task_id].data_provider[lastest_id].registered = true;
		    task[task_id].data_provider[lastest_id].submited = false;
		    task[task_id].data_provider[lastest_id].claimed = false;
		}
		task[task_id].data_provider_id[provider] = lastest_id + 1;
		task[task_id].register_count += 1;
		//emit Register(data_provider, task_id);
		emit Register(task_id,task[task_id].register_count);
		if(task[task_id].register_count == task[task_id].request.target) {
			nextStage(task_id);
			//emit RegisterCollected(task_id);
		}
	}
    
	function submit(uint task_id, bytes calldata data,bytes calldata proof) external {
		require(atStage(task_id, Stages.submit));
		address provider = msg.sender;
		uint id = task[task_id].data_provider_id[provider];
		require (!(id> task[task_id].register_count || id ==0));
		require (task[task_id].data_provider[id-1].registered == true && task[task_id].data_provider[id-1].submited == false);
		
		task[task_id].data_provider[id-1].submited = true;
		task[task_id].data_provider[id-1].submit_data = data;
		task[task_id].data_provider[id-1].submit_proof = proof;
		task[task_id].submit_count += 1;
		//emit Submit(data_provider, task_id);
		emit Submit(task_id, task[task_id].submit_count);
		if(task[task_id].submit_count == task[task_id].register_count){
			nextStage(task_id);
			//emit SubmitCollected(task_id);
		}
	}
	
	// only for demo
	function registerAndSubmit(uint task_id, bytes calldata data,bytes calldata proof) external {
	    require(atStage(task_id, Stages.register));
		address provider = msg.sender;
		uint id = task[task_id].data_provider_id[provider];
		uint lastest_id = task[task_id].register_count;
		require(id == 0 || id > lastest_id || task[task_id].data_provider[id-1].account != provider);
		if (task[task_id].data_provider.length == lastest_id) {
		    task[task_id].data_provider.push(DataProvider(provider,true,false,false,false,"0x1","0x1"));
		} else{
		    task[task_id].data_provider[lastest_id].account = provider;
		    task[task_id].data_provider[lastest_id].registered = true;
		    task[task_id].data_provider[lastest_id].submited = true;
		    task[task_id].data_provider[lastest_id].claimed = false;
		}
		task[task_id].data_provider_id[provider] = lastest_id + 1;
		task[task_id].data_provider[lastest_id].submit_data = data;
		task[task_id].data_provider[lastest_id].submit_proof = proof;
		task[task_id].register_count += 1;
		task[task_id].submit_count += 1;
		//emit Register(data_provider, task_id);
		emit Register(task_id,task[task_id].register_count);
		emit Submit(task_id, task[task_id].submit_count);
		
		if(task[task_id].register_count == task[task_id].request.target) {
		    Stages oldStage = task[task_id].stage;
			task[task_id].stage = Stages.aggregate;
			emit StageTransfer(task_id, uint(task[task_id].stage), uint(oldStage));
			//emit RegisterCollected(task_id);
		}
	}
	
	// only for demo use, force to proceed to aggregate stage 
	function stopRegisterAndSubmit(uint task_id) external {
	    require(atStage(task_id, Stages.register));
	    require(task[task_id].owner == msg.sender);
	    
	    Stages oldStage = task[task_id].stage;
		task[task_id].stage = Stages.aggregate;
		emit StageTransfer(task_id, uint(task[task_id].stage), uint(oldStage));
	}

	function aggregate(uint task_id, bytes calldata aggregation, bytes calldata qualifiedSets,bytes calldata share, bytes calldata attestatino) external {
	    require(atStage(task_id, Stages.aggregate));
	   	require(task[task_id].service_provider.account == msg.sender);
	   	require(qualifiedSets.length == task[task_id].register_count, "wrong qualifiedSets length");
	   	task[task_id].aggregate.aggregation = aggregation;
	   	task[task_id].aggregate.share = share;
	   	task[task_id].aggregate.attestatino = attestatino;
	   	setQualifiedDataProvider(task_id,qualifiedSets);
	   	emit Aggregate(task_id,aggregation);
		nextStage(task_id);
		//emit Aggregate(task_id);
	}

	function approve(uint task_id) external {
	    require(atStage(task_id, Stages.approve));
	    require(task[task_id].owner == msg.sender);
	    emit Approve(task_id);
		nextStage(task_id);
		//emit Approve(task_id);
	}

	function claim(uint task_id) external {
		require(atStage(task_id, Stages.claim));
	    address payable user = msg.sender;
	    uint id = task[task_id].data_provider_id[user];
	    bool is_qualified_data_provider = !(id> task[task_id].register_count || id == 0) &&
	    !task[task_id].data_provider[id-1].claimed &&
	    task[task_id].data_provider[id-1].qualified;
	    bool is_service_provider = ( user == task[task_id].service_provider.account && !task[task_id].service_provider.claimed);
	    require (is_qualified_data_provider || is_service_provider);
		if (is_service_provider) {
			//user.transfer(task[task_id].request.service_fee);  // pay ether
			token.transfer(user,task[task_id].request.service_fee); //pay token
			task[task_id].service_provider.claimed = true;
			task[task_id].claim_count +=1;
		}
		if (is_qualified_data_provider){
			uint reward = task[task_id].request.data_fee/task[task_id].request.target;
// 			address(user).transfer(reward);  // pay ether 
            token.transfer(user,reward);  // pay token
			task[task_id].data_provider[id-1].claimed = true;
			task[task_id].claim_count +=1;
		}
		emit Claim(task_id,task[task_id].claim_count);
		if(task[task_id].claim_count == task[task_id].qualified_count + 1 ) {  // number of data_provider + service_provider
			task[task_id].busy = false;
			task[task_id].register_count = 0;
    		task[task_id].submit_count = 0;
    		task[task_id].claim_count = 0;
    		delete task[task_id].request;
			nextStage(task_id);
		}
	}

}