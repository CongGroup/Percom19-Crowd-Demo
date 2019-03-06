import {AbiCoder} from 'web3-eth-abi';
import {Accounts} from 'web3-eth-accounts';

const config = require('../config/agg.json');
const coder = new AbiCoder();
const accounts = new Accounts('');

function encodeFunction(funcName,...args) {
    let jsonInterface = searchForJsonInterface(funcName);
    // console.log(jsonInterface);
    return coder.encodeFunctionCall(jsonInterface,args)
}

function getBaseTxObject() {
    return {
        to: config.address,
        gasPrice: 0,
        gas:"3000000",
        value:"0",
    }
}

function searchForJsonInterface(funcName) {
    for(let i =0;i<config.abi.length;++i) {
        if (config.abi[i].name === funcName && config.abi[i].type === 'function') {
            return config.abi[i];
        }
    }
}

function signTx(tx,privateKey) {
    return accounts.signTransaction(tx,privateKey).then(signedTx=>{
        return signedTx.rawTransaction;
    })
}

const createAccount = ()=> {
    let wallet =  accounts.create();
    return {
        address: wallet.address,
        privateKey: wallet.privateKey
    }
};

export { createAccount, signTx, encodeFunction,getBaseTxObject};
