package com.exchain.web3.util;

import org.web3j.contracts.eip20.generated.ERC20;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.DefaultGasProvider;

import java.io.IOException;
import java.math.BigInteger;


class Sample {

    static public ERC20 getERC20Instance(Credentials credentials) throws IOException {
        // Connect Web3j to the Blockchain
        String rpcEndpoint = "https://exchaintestrpc.okex.org";
        Web3j web3j = Web3j.build(new HttpService(rpcEndpoint));

        RawTransactionManager txManager = new RawTransactionManager(web3j, credentials, Byte.parseByte(web3j.netVersion().send().getNetVersion()));

        // Load the USDT contract
        String contractAddress = "0xe579156f9decc4134b5e3a30a24ac46bb8b01281";
        return ERC20.load(contractAddress, web3j, txManager, new DefaultGasProvider());
    }

    static public void main(String[] args) throws Exception {
        // prepare private key 1
        Credentials credential1 = Credentials.create("private key1");
        String addr1 = credential1.getAddress();
        System.out.println("address: " + addr1);

        // prepare private key 2
        Credentials credential2 = Credentials.create("private key2");
        String addr2 = credential2.getAddress();
        System.out.println("address: " + addr2);

        ERC20 usdt = getERC20Instance(credential1);

        String symbol = usdt.symbol().send();
        String name = usdt.name().send();
        BigInteger decimal = usdt.decimals().send();

        System.out.println("symbol: " + symbol);
        System.out.println("name: " + name);
        System.out.println("decimal: " + decimal.intValueExact());

        // transfer
        BigInteger balance1 = usdt.balanceOf(addr1).send();
        System.out.println("balance (" + addr1 + ")=" + balance1.toString());

        TransactionReceipt receipt = usdt.transfer(addr2, new BigInteger("25")).send();
        System.out.println("Transaction hash: " + receipt.getTransactionHash());

        balance1 = usdt.balanceOf(addr1).send();
        System.out.println("balance (" + addr1 + ")=" + balance1.toString());

        BigInteger balance2 = usdt.balanceOf(addr2).send();
        System.out.println("balance (" + addr2 + ")=" + balance2.toString());

        // approve
        TransactionReceipt receipt1 = usdt.approve(addr2, new BigInteger("100000")).send();
        System.out.println("Transaction hash: " + receipt1.getTransactionHash());

        BigInteger approveAmt = usdt.allowance(addr1, addr2).send();
        System.out.println("allowance (" + addr2 + ")=" + approveAmt.toString());

        // transferFrom
        usdt = getERC20Instance(credential2);
        TransactionReceipt receipt2 = usdt.transferFrom(addr1, "0x0db6b797e64666d4b36b13e5dc6fcd4661893ac8", new BigInteger("50000")).send();
        System.out.println("Transaction hash: " + receipt2.getTransactionHash());

        approveAmt = usdt.allowance(addr1, addr2).send();
        System.out.println("allowance (" + addr2 + ")=" + approveAmt.toString());
    }
}
