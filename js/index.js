const HDWalletProvider = require('@truffle/hdwallet-provider');
const Web3 = require('web3');

const privkey = "189021ea03a74055cd974d4181e1f7526f1a87c8be4b318407a14385606b10df";
const address = "0xd01bf1F0dB0E0F9998Ec01a45Cfc03116D0224bE";

const provider = new HDWalletProvider(privkey, "https://exchaintestrpc.okex.org/");
const web3 = new Web3(provider);

const contractAbi = [{"constant":false,"inputs":[{"internalType":"uint256","name":"delta","type":"uint256"}],"name":"add","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"subtract","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"counter","type":"uint256"}],"name":"Changed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"counter","type":"uint256"}],"name":"Added","type":"event"}]
const contractAddress = '0x79BE5cc37B7e17594028BbF5d43875FDbed417db'
const myContract = new web3.eth.Contract(contractAbi, contractAddress);


myContract.methods.getCounter().call((err, result) => {
    console.log(result)
})



async function doWrite() {
    const nounce = await web3.eth.getTransactionCount(address);

    const addTx = {
        // this could be provider.addresses[0] if it exists
        from: address,
        // target address, this could be a smart contract address
        to: contractAddress,
        // optional if you want to specify the gas limit
        gas: 2000000,

        nonce: nounce,

        // this encodes the ABI of the method and the arguements
        data: myContract.methods.add(10).encodeABI()
    };

    const addSignPromise = web3.eth.accounts.signTransaction(addTx, privkey);

    addSignPromise.then((signedTx) => {

        // raw transaction string may be available in .raw or
        // .rawTransaction depending on which signTransaction
        // function was called
        const sentTx = web3.eth.sendSignedTransaction(signedTx.raw || signedTx.rawTransaction);

        sentTx.on("receipt", receipt => {
            // do something when receipt comes back
            console.log(receipt)

            myContract.methods.getCounter().call((err, result) => {
                console.log('add', result, 'nounce', nounce)
            })

            doSubtract();
        });

        sentTx.on("error", err => {
            // do something on transaction error
            console.log(err)
        });

    }).catch((err) => {

        // do something when promise fails

    });

    function doSubtract () {
        const subtractTx = {
            // this could be provider.addresses[0] if it exists
            from: address,
            // target address, this could be a smart contract address
            to: contractAddress,
            // optional if you want to specify the gas limit
            gas: 2000000,

            nonce: (nounce + 1),

            // this encodes the ABI of the method and the arguements
            data: myContract.methods.subtract().encodeABI()
        };

        const subtractSignPromise = web3.eth.accounts.signTransaction(subtractTx, privkey);

        subtractSignPromise.then((signedTx) => {

            // raw transaction string may be available in .raw or
            // .rawTransaction depending on which signTransaction
            // function was called
            const sentTx = web3.eth.sendSignedTransaction(signedTx.raw || signedTx.rawTransaction);

            sentTx.on("receipt", receipt => {
                // do something when receipt comes back
                console.log(receipt)

                myContract.methods.getCounter().call((err, result) => {
                    console.log('subtract', result, 'nounce', nounce)
                })

                // myContract.getPastEvents('MyEvent', {}, function (error, event) {
                //     console.log('MyEvent------------------------', event);
                // }).on('Added', function (event) {
                //     console.log('Added', event);
                // }).on('Changed', function (event) {
                //     console.log('Changed', event);
                // })
            });

            sentTx.on("error", err => {
                // do something on transaction error
                console.log(err)
            });

        }).catch((err) => {

            // do something when promise fails

        });
    }





}


doWrite()





