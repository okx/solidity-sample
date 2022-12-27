const Web3 = require('web3');

async function main() {

//contract ABI
const systemContractABI = [
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_data",
          "type": "string"
        }
      ],
      "name": "invoke",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]
const stakingMSGHelperABI = [
    {
      "inputs": [
        {
          "internalType": "string[]",
          "name": "_validatorAddresses",
          "type": "string[]"
        }
      ],
      "name": "genAddSharesMsg",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_amount",
          "type": "string"
        }
      ],
      "name": "genDepositMsg",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "genWithdrawAllRewardsMsg",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_amount",
          "type": "string"
        }
      ],
      "name": "genWithdrawMsg",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_str",
          "type": "string"
        }
      ],
      "name": "stringToHexString",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "pure",
      "type": "function"
    }
  ]

//user privateKey, add you private key
const userPrivateKey = ""
//user address, add you address
const userAddress = ""
//validator address
const validatorAddresses = ""
const validatorAddresses2 = ""
//the OKT you want deposit,Minimum amount is 1.(1 = 1 OKT)
const amount = "100"

//contract address in okc main net
const systemContractAddress = "0xd6bce454316b8ddFb76bB7bb1B57B8942B09Acd5"
//stakingMSGHelper proxy contract address in okc main net
const stakingMSGHelperContractAddress = "0x1b29c875Bd7Ec9a12C29fc6eeF8E451207352EF3"

//init web3, connect okc main net
const web3 = new Web3(new Web3.providers.HttpProvider("https://exchainrpc.okex.org"));

//init contract object
const systemContract = new web3.eth.Contract(systemContractABI, systemContractAddress);
const stakingMSGHelperContract = new web3.eth.Contract(stakingMSGHelperABI, stakingMSGHelperContractAddress);

//deposit
//get data
let backData = await stakingMSGHelperContract.methods.genDepositMsg(amount).call({from: userAddress})

//encode
let encodeData  = await systemContract.methods.invoke(backData).encodeABI();

//sign
let sign = await web3.eth.accounts.signTransaction({
    gas: 500000,
    to: systemContractAddress,
    data: encodeData,
}, userPrivateKey)

//send tx
result = await web3.eth.sendSignedTransaction(sign.rawTransaction)

//add validator addresses
//get data
backData = await stakingMSGHelperContract.methods.genAddSharesMsg([validatorAddresses, validatorAddresses2]).call({from: userAddress})

//encode
encodeData  = await systemContract.methods.invoke(backData).encodeABI();

//sign
sign = await web3.eth.accounts.signTransaction({
    gas: 500000,
    to: systemContractAddress,
    data: encodeData,
}, userPrivateKey)

//send tx
result = await web3.eth.sendSignedTransaction(sign.rawTransaction)

//withdraw all rewards
//get data
backData = await stakingMSGHelperContract.methods.genWithdrawAllRewardsMsg().call({from: userAddress})

//encode
encodeData  = await systemContract.methods.invoke(backData).encodeABI();

//sign
sign = await web3.eth.accounts.signTransaction({
    gas: 500000,
    to: systemContractAddress,
    data: encodeData,
}, userPrivateKey)

//send tx
result = await web3.eth.sendSignedTransaction(sign.rawTransaction)

//withdraw
//get data
backData = await stakingMSGHelperContract.methods.genWithdrawMsg(amount).call({from: userAddress})

//encode
encodeData  = await systemContract.methods.invoke(backData).encodeABI();

//sign
sign = await web3.eth.accounts.signTransaction({
    gas: 500000,
    to: systemContractAddress,
    data: encodeData,
}, userPrivateKey)

//send tx
result = await web3.eth.sendSignedTransaction(sign.rawTransaction)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
  