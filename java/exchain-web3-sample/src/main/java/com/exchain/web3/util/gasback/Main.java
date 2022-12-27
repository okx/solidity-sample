package com.exchain.web3.util.gasback;

import com.exchain.web3.util.ContractUtil;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.*;
import org.web3j.abi.datatypes.generated.StaticArray1;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.RawTransaction;
import org.web3j.crypto.TransactionEncoder;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import org.web3j.protocol.http.HttpService;
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

    // define testnet systemContractAddress. mainnet systemContractAddress is here: https://www.oklink.com/zh-cn/okc/address/0xe9196e65a0b6705777fbe829dfa94ec8b9f2ba48
    private static String systemContractAddress = "0x7e5E6AF6424BE7A835313777Fc6E0d1912e52Fc8";
    // define testnet gasbackHelplerContractAddress. mainnet gasbackHelplerContractAddress is here: https://www.oklink.com/zh-cn/okc/address/0x0dd08b74c111d148751f38f02ab0c3408ead7d18
    private static String gasbackHelplerContractAddress = "0x56CeE7c20F4C83996e72b65a7b410fcE4C3b52B0";
    private static String targerContractAddress = "0x7C9eBC51F32ED4760AEADeF365875ef5727442Bc";
    private static String privateKey = "";
    private static String fromAddress = "0xbbE4733d85bc2b90682147779DA49caB38C0aA1F";

    private static Web3j web3j = Web3j.build(new HttpService(url));

    public static void main(String[] args) throws TransactionException, IOException, ExecutionException, InterruptedException {

        // Deposit
        List<Type> params = new ArrayList<>();

        params.add(new Address(targerContractAddress));
        params.add(new Address(fromAddress));
        params.add(new DynamicArray(new Uint256(394658)));

        TransactionReceipt receipt = invoke("genRegisterMsg", params);
        System.out.println(receipt);

        params.clear();
        params.add(new Address(targerContractAddress));
        params.add(new Address("0x04A987fa1Bd4b2B908e9A3Ca058cc8BD43035991"));
        receipt = invoke("genUpdateMsg", params);
        System.out.println(receipt);

        params.clear();
        params.add(new Address(targerContractAddress));
        receipt = invoke("genCancelMsg", params);
        System.out.println(receipt);
    }


    public static TransactionReceipt invoke(String name, List<Type> args) throws TransactionException, IOException, ExecutionException, InterruptedException {
        List<TypeReference<?>> outputParameters = new ArrayList<>();
        TypeReference<Utf8String> typeReference = new TypeReference<Utf8String>() {
        };
        outputParameters.add(typeReference);

        // gen msg
        List<Type> ret = ContractUtil.ethCall(web3j, fromAddress, gasbackHelplerContractAddress, name, args, outputParameters);

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
        byte[] signMessage = TransactionEncoder.signMessage(rawTransaction, Byte.valueOf("65"), credentials);

        String hash = web3j.ethSendRawTransaction(Numeric.toHexString(signMessage)).sendAsync().get().getTransactionHash();

        System.out.println("hash:" + hash);
        return receiptProcessor.waitForTransactionReceipt(hash);
    }
}
