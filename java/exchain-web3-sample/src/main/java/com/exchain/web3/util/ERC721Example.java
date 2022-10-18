package com.exchain.web3.util;

import org.web3j.contracts.eip20.generated.ERC20;
import org.web3j.contracts.eip721.generated.ERC721;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.ContractGasProvider;
import org.web3j.tx.gas.DefaultGasProvider;

import java.io.IOException;
import java.math.BigInteger;
import java.nio.charset.Charset;

public class ERC721Example {

    public static ERC721 getERC721Instance(Credentials credentials) throws IOException {
        // Connect Web3j to the Blockchain
        String rpcEndpoint = "https://exchaintestrpc.okex.org";
        Web3j web3j = Web3j.build(new HttpService(rpcEndpoint));

        RawTransactionManager txManager = new RawTransactionManager(web3j, credentials, Byte.parseByte(web3j.netVersion().send().getNetVersion()));

        // Load the my art contract
        String contractAddress = "0x2366E30fC4Dd41FEc5c897087b2d3746a2Ed47cc";
        return ERC721.load(contractAddress, web3j, txManager, new DefaultGasProvider());
    }

    public static void main(String[] args) throws Exception {


        // prepare private key 1
        Credentials credential1 = Credentials.create("private key1");
        String addr1 = credential1.getAddress();
        System.out.println("address1: " + addr1);

        // prepare private key 2
        Credentials credential2 = Credentials.create("private key2");
        String addr2 = credential2.getAddress();
        System.out.println("address2: " + addr2);

        ERC721 fas = getERC721Instance(credential1);

        // query balance
        BigInteger balance1 = fas.balanceOf(addr1).send();
        System.out.println(addr1 + "balance: " + balance1);


        // query owner
        // 0xbbe4733d85bc2b90682147779da49cab38c0aa1f
        String ownerOf = fas.ownerOf(BigInteger.ZERO).send();
        System.out.println("owner of 0 is:" + ownerOf);

//        // transferFrom
        TransactionReceipt receipt = fas.transferFrom(addr1, addr2, BigInteger.valueOf(4), BigInteger.valueOf(0)).send();
        System.out.println(receipt);

        // safeTransferFrom
        receipt = fas.safeTransferFrom(addr1, addr2, BigInteger.valueOf(10), BigInteger.valueOf(0)).send();
        System.out.println(receipt);

        // safeTransferFrom
        receipt = fas.safeTransferFrom(addr1, addr2, BigInteger.valueOf(11), "hello".getBytes(Charset.forName("utf-8")), BigInteger.valueOf(0)).send();
        System.out.println(receipt);

        // query addr2 balance
        BigInteger balance2 = fas.balanceOf(addr2).send();
        System.out.println(addr2 + "balance: " + balance2);

        // approve 12 to addr2
        receipt = fas.approve(addr2, BigInteger.valueOf(12), BigInteger.valueOf(0)).send();
        System.out.println(receipt);

        // query Approver of 12
        String approver = fas.getApproved(BigInteger.valueOf(12)).send();
        System.out.println("the approver of 12 is:" + approver);


        // setApprovalForAll
        receipt = fas.setApprovalForAll(addr2, true).send();
        System.out.println(receipt);

        boolean yes = fas.isApprovedForAll(addr1, addr2).send();
        System.out.println(addr2 + "is operator :" + yes);


        fas = getERC721Instance(credential2);
        String addr3 = "0x9ed8C8667cBF24bc01b427AD373774034A83Dd88";
        // addr2 operator the transfer: transfer 12 from addr1 to add3
        receipt = fas.transferFrom(addr1, addr3, BigInteger.valueOf(12), BigInteger.ZERO).send();
        System.out.println(receipt);
    }
}
