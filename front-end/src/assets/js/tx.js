// import {AbiCoder} from 'web3-eth-abi';
import {Accounts} from 'web3-eth-accounts';
// const coder = new AbiCoder();
const accounts = new Accounts('ws://127.0.0.1:8650');

// function add0x (input) {
//     if (typeof(input) !== 'string') {
//         return input;
//     }
//     if (input.length < 2 || input.slice(0,2) !== '0x') {
//         return '0x' + input;
//     }
//     return input;
// }

// function assembleFunction (jsonInterface, args) {
//     // txObject contains gasPrice, gasLimit, nonce, to, value
//     // let jsonInterface = {
//     //     "name": "test",
//     //     "type": "function",
//     //     "inputs": [
//     //         {
//     //             "name": "c",
//     //             "type": "int256"
//     //         },
//     //         {
//     //             "name": "b",
//     //             "type": "int256"
//     //         }
//     //     ],
//     // };
//
//     let txData = coder.encodeFunctionCall(jsonInterface, args);
//
//     return txData
// }

// function signTx(tx, privateKey) {
//     return accounts.signTransaction(tx,privateKey)
// }

const createAccount = ()=> {
    let wallet =  accounts.create();
    return {
        address: wallet.address,
        privateKey: wallet.privateKey
    }
};

export { createAccount};
