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
const gasBackContractABI = [
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_contract",
          "type": "address"
        }
      ],
      "name": "genCancelMsg",
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
          "internalType": "address",
          "name": "_contract",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_recipientAddress",
          "type": "address"
        },
        {
          "internalType": "uint256[]",
          "name": "_nonces",
          "type": "uint256[]"
        }
      ],
      "name": "genRegisterMsg",
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
          "internalType": "address",
          "name": "_contract",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_recipientAddress",
          "type": "address"
        }
      ],
      "name": "genUpdateMsg",
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

//user privateKey, add your private key
const userPrivateKey = ""
//user address, add your address
const userAddress = ""
//address to receive rewards
const recipientAddress = ""
//address that the new receiver address
const updateRecipientAddress = ""
//contract address who want register
const contractAddressNeedRegister = ""
// nonce when deploy contract
// eg:const nonce = 81
const nonce = 

//contract address in okc main net
const systemContractAddress = "0xd6bce454316b8ddFb76bB7bb1B57B8942B09Acd5"
//gasBackMSGHelper proxy contract address in okc main net
const gasBackContractAddress = "0x0DD08B74c111D148751f38f02ab0C3408ead7d18"

//init web3, connect okc main net
const web3 = new Web3(new Web3.providers.HttpProvider("https://exchainrpc.okex.org"));

//init contract object
const systemContract = new web3.eth.Contract(systemContractABI, systemContractAddress);
const gasBackContract = new web3.eth.Contract(gasBackContractABI, gasBackContractAddress);

//register
//get data
let backData = await gasBackContract.methods.genRegisterMsg(contractAddressNeedRegister,recipientAddress, [nonce,]).call({from: userAddress})

//encode
let encodeData  = await systemContract.methods.invoke(backData).encodeABI();

//sign
let sign = await web3.eth.accounts.signTransaction({
    gas: 500000,
    to: systemContractAddress,
    data: encodeData
}, userPrivateKey)

//send tx
result = await web3.eth.sendSignedTransaction(sign.rawTransaction)


//update
//get data
backData = await gasBackContract.methods.genUpdateMsg(contractAddressNeedRegister, updateRecipientAddress).call({from: userAddress})

//encode
encodeData  = await systemContract.methods.invoke(backData).encodeABI();

//sign
sign = await web3.eth.accounts.signTransaction({
    gas: 500000,
    to: systemContractAddress,
    data: encodeData
}, userPrivateKey)

//send tx
result = await web3.eth.sendSignedTransaction(sign.rawTransaction)


//cancel
//get data
backData = await gasBackContract.methods.genCancelMsg(contractAddressNeedRegister).call({from: userAddress})

//encode
encodeData  = await systemContract.methods.invoke(backData).encodeABI();

//sign
sign = await web3.eth.accounts.signTransaction({
    gas: 500000,
    to: systemContractAddress,
    data: encodeData
}, userPrivateKey)

//send tx
result = await web3.eth.sendSignedTransaction(sign.rawTransaction)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });