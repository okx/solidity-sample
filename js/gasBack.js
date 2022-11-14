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
          "name": "_withdrawerAddress",
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
          "name": "_withdrawerAddress",
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
//contract address who want register
const contractAddressNeedRegister = ""
// nonce when deploy contarct
//eg:const nonce = 81
const nonce = 

//contract address in okc test net
//note: okc main net is: 0xd6bce454316b8ddFb76bB7bb1B57B8942B09Acd5
const systemContarctAddress = "0x727d14EfC4FB5281A18A6d62BCf486a1cF4d2210"
//gasBackMSGHelper proxy contract address in okc test net
//note: okc main net is: 0x0DD08B74c111D148751f38f02ab0C3408ead7d18
const gasBackContractAddress = "0x9e472f77e2A5C8f09B237273960c776ddE1D98C1"

//init web3, connect okc test net
const web3 = new Web3(new Web3.providers.HttpProvider("https://exchaintestrpc.okex.org"));

//init contract object
const systemContract = new web3.eth.Contract(systemContractABI, systemContarctAddress);
const gasBackContract = new web3.eth.Contract(gasBackContractABI, gasBackContractAddress);

//get data
let backData = await gasBackContract.methods.genRegisterMsg(contractAddressNeedRegister,userAddress, [nonce,]).call({from: userAddress})

//encode
let encodeData  = await systemContract.methods.invoke(backData).encodeABI();

//sign
let sign = await web3.eth.accounts.signTransaction({
    gas: 500000,
    to: systemContarctAddress,
    data: encodeData
}, userPrivateKey)

//send tx
result = await web3.eth.sendSignedTransaction(sign.rawTransaction)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });