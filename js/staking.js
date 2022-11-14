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
//the OKT you want deposit
const amount = "100"

//contract address in okc test net
//note: okc main net is: 0xd6bce454316b8ddFb76bB7bb1B57B8942B09Acd5
const systemContarctAddress = "0x727d14EfC4FB5281A18A6d62BCf486a1cF4d2210"
//stakingMSGHelper proxy contract address in okc test net
//note: okc main net is: 0x1b29c875Bd7Ec9a12C29fc6eeF8E451207352EF3
const stakingMSGHelperContractAddress = "0x37AFc0Fe2b14acC1593C4435F3B53ac0335Ee292"

//init web3, connect okc test net
const web3 = new Web3(new Web3.providers.HttpProvider("https://exchaintestrpc.okex.org"));

//init contract object
const systemContract = new web3.eth.Contract(systemContractABI, systemContarctAddress);
const stakingMSGHelperContract = new web3.eth.Contract(stakingMSGHelperABI, stakingMSGHelperContractAddress);

//get data
let backData = await stakingMSGHelperContract.methods.genDepositMsg(amount).call({from: userAddress})

//encode
let encodeData  = await systemContract.methods.invoke(backData).encodeABI();

//sign
let sign = await web3.eth.accounts.signTransaction({
    gas: 500000,
    to: systemContarctAddress,
    data: encodeData,
}, userPrivateKey)

//send tx
result = await web3.eth.sendSignedTransaction(sign.rawTransaction)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
  