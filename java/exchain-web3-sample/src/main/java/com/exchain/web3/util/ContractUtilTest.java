package com.exchain.web3.util;

import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Uint;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.DefaultGasProvider;
import org.web3j.tx.response.PollingTransactionReceiptProcessor;
import org.web3j.tx.response.TransactionReceiptProcessor;

import java.math.BigInteger;
import java.util.*;

/**
 * @author shaoyun.zhan
 * @date 2021/5/17
 * <p>
 * 描述：
 */
public class ContractUtilTest {

    private static Web3j web3j;

    private static String url = "https://exchaintestrpc.okex.org";

    private static String fromAddress = "0x04A987fa1Bd4b2B908e9A3Ca058cc8BD43035991";

    private static String contractAddress = "0x79BE5cc37B7e17594028BbF5d43875FDbed417db";

    private static String privateKey = "89c81c304704e9890025a5a91898802294658d6e4034a11c6116f4b129ea12d3";

    public static void main(String[] args) throws Exception {

        web3j = Web3j.build(new HttpService(url));

        //operate
        System.out.println(add(BigInteger.valueOf(1)));
        System.out.println(subtract());
        System.out.println(getCounter(web3j, fromAddress, contractAddress));

        web3j.shutdown();
    }

    /**
     * getCounter
     * read contract with params
     */
    public static BigInteger getCounter(Web3j web3j, String fromAddress, String contractAddress) {
        String methodName = "getCounter";

        List<Type> inputParameters = new ArrayList<>();

        List<TypeReference<?>> outputParameters = new ArrayList<>();
        TypeReference<Uint256> typeReference = new TypeReference<Uint256>() {
        };
        outputParameters.add(typeReference);


        List<Type> results = ContractUtil.ethCall(web3j, fromAddress, contractAddress, methodName, inputParameters, outputParameters);
        BigInteger counter = (BigInteger) results.get(0).getValue();
        return counter;
    }


    /**
     * add
     *
     * @param
     * @return
     * @throws Exception
     */
    private static String add(BigInteger addNum) throws Exception {

        // Load an account
        Credentials credentials = Credentials.create(privateKey);

        // Contract and functions
        Function function = new Function("add", // Function name
                Arrays.asList(new Uint(addNum)), // Function input parameters
                Collections.emptyList()); // Function returned parameters

        //Encode function values in transaction data format
        String txData = FunctionEncoder.encode(function);

        // RawTransactionManager use a wallet (credential) to create and sign transaction
        TransactionManager txManager = new RawTransactionManager(web3j, credentials);

        // Send transaction
        String txHash = txManager.sendTransaction(
                DefaultGasProvider.GAS_PRICE,
                DefaultGasProvider.GAS_LIMIT,
                contractAddress,
                txData,
                BigInteger.ZERO).getTransactionHash();

        // Wait for transaction to be mined
        TransactionReceiptProcessor receiptProcessor = new PollingTransactionReceiptProcessor(
                web3j,
                TransactionManager.DEFAULT_POLLING_FREQUENCY,
                TransactionManager.DEFAULT_POLLING_ATTEMPTS_PER_TX_HASH);
        TransactionReceipt txReceipt = receiptProcessor.waitForTransactionReceipt(txHash);
        return txReceipt.getTransactionHash();
    }

    /**
     * subtract
     *
     * @param
     * @return
     * @throws Exception
     */
    private static String subtract() throws Exception {

        // Load an account
        Credentials credentials = Credentials.create(privateKey);

        // Contract and functions
        Function function = new Function("subtract", // Function name
                new LinkedList<>(), // Function input parameters
                Collections.emptyList()); // Function returned parameters

        //Encode function values in transaction data format
        String txData = FunctionEncoder.encode(function);

        // RawTransactionManager use a wallet (credential) to create and sign transaction
        TransactionManager txManager = new RawTransactionManager(web3j, credentials);

        // Send transaction
        String txHash = txManager.sendTransaction(
                DefaultGasProvider.GAS_PRICE,
                DefaultGasProvider.GAS_LIMIT,
                contractAddress,
                txData,
                BigInteger.ZERO).getTransactionHash();

        // Wait for transaction to be mined
        TransactionReceiptProcessor receiptProcessor = new PollingTransactionReceiptProcessor(
                web3j,
                TransactionManager.DEFAULT_POLLING_FREQUENCY,
                TransactionManager.DEFAULT_POLLING_ATTEMPTS_PER_TX_HASH);
        TransactionReceipt txReceipt = receiptProcessor.waitForTransactionReceipt(txHash);
        return txReceipt.getTransactionHash();
    }

}
