package com.exchain.web3.util.staking;

import com.exchain.web3.util.ContractUtil;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.DynamicArray;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.StaticArray1;
import org.web3j.abi.datatypes.generated.Uint64;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.RawTransaction;
import org.web3j.crypto.TransactionEncoder;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.DefaultGasProvider;
import org.web3j.tx.response.PollingTransactionReceiptProcessor;
import org.web3j.tx.response.TransactionReceiptProcessor;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;

import java.io.IOException;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.ExecutionException;

public class Main {

    private static String url = "https://exchaintestrpc.okex.org";

    private static String systemContractAddress = "0x7e5E6AF6424BE7A835313777Fc6E0d1912e52Fc8";
    private static String stakingContractAddress = "0xE6947172736e3ed558Ec375cf16C62488627a18e";

    private static String privateKey = "89c81c304704e9890025a5a91898802294658d6e4034a11c6116f4b129ea12d3";
    private static String fromAddress = "0x04A987fa1Bd4b2B908e9A3Ca058cc8BD43035991";

    private static Web3j web3j = Web3j.build(new HttpService(url));

    public static void main(String[] args) throws TransactionException, IOException, ExecutionException, InterruptedException {

        // Deposit
        List<Type> params = new ArrayList<>();
        params.add(new Utf8String("100"));

        TransactionReceipt receipt = invoke("genDepositMsg", params);
        System.out.println(receipt);

        params.clear();
        params.add(new Utf8String("2"));
        receipt = invoke("genWithdrawMsg", params);
        System.out.println(receipt);

        params.clear();
        params.add(new DynamicArray(new Utf8String("exvaloper15w73dl3n7qmrk8ssax5w2k45qd29jlkejyc23q")));
        receipt = invoke("genAddSharesMsg", params);
        System.out.println(receipt);

        params.clear();
        receipt = invoke("genWithdrawAllRewardsMsg", params);
        System.out.println(receipt);
    }


    public static TransactionReceipt invoke(String name, List<Type> args) throws TransactionException, IOException, ExecutionException, InterruptedException {
        List<TypeReference<?>> outputParameters = new ArrayList<>();
        TypeReference<Utf8String> typeReference = new TypeReference<Utf8String>() {
        };
        outputParameters.add(typeReference);

        // gen msg
        List<Type> ret = ContractUtil.ethCall(web3j, fromAddress, stakingContractAddress, name, args, outputParameters);

        // invoke systemcontract
        // Load an account
        Credentials credentials = Credentials.create(privateKey);
        // Contract and functions

        Function function = new Function("invoke", ret, Collections.emptyList());

        //Encode function values in transaction data format
        String txData = FunctionEncoder.encode(function);

        TransactionReceiptProcessor receiptProcessor = new PollingTransactionReceiptProcessor(web3j, TransactionManager.DEFAULT_POLLING_FREQUENCY, TransactionManager.DEFAULT_POLLING_ATTEMPTS_PER_TX_HASH);
        BigInteger amountWei = Convert.toWei("0.00000005", Convert.Unit.ETHER).toBigInteger();
        System.out.println(amountWei);

        //sign tx
        BigInteger nonce = web3j.ethGetTransactionCount(fromAddress, DefaultBlockParameterName.PENDING).send().getTransactionCount();
        RawTransaction rawTransaction = RawTransaction.createTransaction(nonce, DefaultGasProvider.GAS_PRICE, DefaultGasProvider.GAS_LIMIT, systemContractAddress, amountWei, txData);
        System.out.println(credentials.getAddress());
        byte[] signMessage = TransactionEncoder.signMessage(rawTransaction,Byte.valueOf("65"), credentials);

        String hash = web3j.ethSendRawTransaction(Numeric.toHexString(signMessage)).sendAsync().get().getTransactionHash();

        System.out.println("hash:" + hash);
        return receiptProcessor.waitForTransactionReceipt(hash);
    }
}
