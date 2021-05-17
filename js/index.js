const Web3 = require('web3')
const rpcURL = 'https://exchaintestrpc.okex.org' // Your RCP URL goes here
const web3 = new Web3(rpcURL)

const abi = [{"constant":false,"inputs":[{"internalType":"uint256","name":"delta","type":"uint256"}],"name":"add","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"subtract","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"counter","type":"uint256"}],"name":"Changed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"counter","type":"uint256"}],"name":"Added","type":"event"}]
const address = '0x79BE5cc37B7e17594028BbF5d43875FDbed417db'

const contract = new web3.eth.Contract(abi, address)

contract.methods.getCounter().call((err, result) => {
    console.log(result)
})

// contract.Changed((err, result) => {
//     console.log(result)
// })
// contract.Added((err, result) => {
//     console.log(result)
// })

contract.methods.add(10).call((err, result) => {
    console.log(result)
})
// contract.methods.getCounter().call((err, result) => {
//     console.log(result)
// })

contract.methods.subtract().call((err, result) => {
    console.log(result)
})
// contract.methods.getCounter().call((err, result) => {
//     console.log(result)
// })