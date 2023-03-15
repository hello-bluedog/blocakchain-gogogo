
package com.dataSharing.service;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.PrivateKey;
import java.util.EnumSet;
import java.util.Properties;
import java.util.Set;
import java.util.concurrent.TimeoutException;

import org.hyperledger.fabric.gateway.Contract;
import org.hyperledger.fabric.gateway.ContractException;
import org.hyperledger.fabric.gateway.Gateway;
import org.hyperledger.fabric.gateway.Identities;
import org.hyperledger.fabric.gateway.Identity;
import org.hyperledger.fabric.gateway.Network;
import org.hyperledger.fabric.gateway.Wallet;
import org.hyperledger.fabric.gateway.Wallets;
import org.hyperledger.fabric.gateway.X509Identity;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.User;
import org.hyperledger.fabric.sdk.security.CryptoSuite;
import org.hyperledger.fabric.sdk.security.CryptoSuiteFactory;
import org.hyperledger.fabric_ca.sdk.EnrollmentRequest;
import org.hyperledger.fabric_ca.sdk.HFCAClient;
import org.hyperledger.fabric_ca.sdk.RegistrationRequest;
import org.springframework.stereotype.Service;

import com.dataSharing.entity.ledgerEntity.CoinNumAndCredit;

public class IAService {
	//提供一个登录时认证和交易前认证接口。在每次登录时和每次交易时我都会调用这个验证函数
	private String walletPath = "wallet";
	private String adminName = "admin";
	private String caPemFile = "cert.pem";
	private String adminPW = "adminpw";
	private String connectionPath = "connection.json";
	private static String contractName = "basic";
	public static Network network;
	public IAService() throws Exception {
		verify(String username)；
	}
	public boolean verify(String username) throws Exception {
	   // 获取的钱包
            Wallet wallet = Wallets.newFileSystemWallet(Paths.get("wallet"));
	   // 检查是否注册过
            /*if (wallet.get(username) != null) {
		    System.out.println("Authentication is successful");
		    return true;
            }
	    else {
		    System.out.println("Authentication is unsuccessful");
		    return false;
            }*/
	    //检查证书
	     X509Identity userIdentity = (X509Identity)wallet.get("username");
	     if (userIdentity == null) {
		    System.out.println("Authentication is unsuccessful");
		    return false;
	     }
	    else {
		    System.out.println("Authentication is successful");
		    return true;}
	  }
	
}
