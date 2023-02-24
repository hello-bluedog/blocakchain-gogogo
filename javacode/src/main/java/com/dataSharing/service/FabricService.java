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

@Service
public class FabricService {
	private String walletPath = "wallet";
	private String adminName = "admin";
	private String caPemFile = "cert.pem";
	private String adminPW = "adminpw";
	private String connectionPath = "connection.json";
	private static String contractName = "basic";
	public static Network network;


	public FabricService() throws Exception {
		createAdmin();
		login();
	}

	public Network login() throws IOException {
		Wallet wallet = Wallets.newFileSystemWallet(Paths.get(this.walletPath));
	    Path networkConfigPath = Paths.get(this.connectionPath);
	    Gateway.Builder builder = Gateway.createBuilder();
	    builder.identity(wallet, this.adminName).networkConfig(networkConfigPath).discovery(true);
	    Network network =  builder.connect().getNetwork("mychannel");
	    this.network = network;
	    return network;
	}

	private HFCAClient getCaclient() throws Exception {
		// Create a CA client for interacting with the CA.
	    Properties props = new Properties();
	    props.put("pemFile",
	                this.caPemFile);
	    props.put("allowAllHostNames", "true");
	    HFCAClient caClient = HFCAClient.createNewInstance("https://182.92.193.27:7054", props);
	    CryptoSuite cryptoSuite = CryptoSuiteFactory.getDefault().getCryptoSuite();
	    caClient.setCryptoSuite(cryptoSuite);
	    return caClient;
	}

	private boolean createAdmin() throws Exception {
	    // 创建CA客户端实例
	    HFCAClient caClient = getCaclient();

	    // 获取管理身份的钱包，没有则创建，这里的写法会创建在项目所在目录下，也可以写成绝对路径形式
	    Wallet wallet = Wallets.newFileSystemWallet(Paths.get("wallet"));
	    // 检查是否已经注册过了管理员身份的用户，是的话直接退出
	    if (wallet.get("admin") != null) {
	   		System.out.println("An identity for the admin user \"admin\" already exists in the wallet");
	        return false;
		}

	    // Enroll the admin user, and import the new identity into the wallet.
	    final EnrollmentRequest enrollmentRequestTLS = new EnrollmentRequest();
	    enrollmentRequestTLS.addHost("localhost");
	    enrollmentRequestTLS.setProfile("tls");
	    // 进行注册，得到注册结果
	    Enrollment enrollment = caClient.enroll(this.adminName, this.adminPW, enrollmentRequestTLS);
	    // 利用注册结果生成新的证书
	    Identity user = Identities.newX509Identity("Org1MSP", enrollment);
	    // 把证书加入钱包中
	    wallet.put("admin", user);
	    System.out.println("Successfully enrolled user \"admin\" and imported it into the wallet");
	    return true;
	}


	public boolean createUser(String username) throws Exception {

		    HFCAClient caClient = getCaclient();
		    // 获取的钱包
		    Wallet wallet = Wallets.newFileSystemWallet(Paths.get("wallet"));

		    // 检查是否已经注册过
		    if (wallet.get(username) != null) {
		    	System.out.println("An identity for the user \"appUser\" already exists in the wallet");
		        return false;
		    }
		    // 获取管理员证书
		    X509Identity adminIdentity = (X509Identity)wallet.get("admin");
		    if (adminIdentity == null) {
		    	System.out.println("\"admin\" needs to be enrolled and added to the wallet first");
		        return false;
			}
		    // 根据管理员的信息创建用户实体
		    User admin = new User() {
				@Override
		        public String getName() {
		        	return "admin";
				}
		        @Override
		        public Set<String> getRoles() {
		        	return null;
		        }
				@Override
		        public String getAccount() {
		        	return null;
				}
		        @Override
		        public String getAffiliation() {
		        	return "";
				}
		        @Override
		        public Enrollment getEnrollment() {
		        	return new Enrollment() {
						@Override
		                public PrivateKey getKey() {
		                	return adminIdentity.getPrivateKey();
						}
		                @Override
		                public String getCert() {
		                	return Identities.toPemString(adminIdentity.getCertificate());
						}
					};
				}
		        @Override
		        public String getMspId() {
		        	return "Org1MSP";
				}
			};
		    // Register the user, enroll the user, and import the new identity into the wallet.
		    RegistrationRequest registrationRequest = new RegistrationRequest(username);
		    registrationRequest.setAffiliation("org1.department1");
		    registrationRequest.setEnrollmentID(username);
		    String enrollmentSecret = caClient.register(registrationRequest, admin);
		    System.out.println(enrollmentSecret);
		    Enrollment enrollment = caClient.enroll(username, enrollmentSecret);
		    X509Identity user = Identities.newX509Identity("Org1MSP", enrollment);
		    wallet.put(username, user);
		    System.out.println("Successfully enrolled user \"appUser\" and imported it into the wallet");
		    return true;
		}
	public static String AddNewUserItem(CoinNumAndCredit newUser) throws ContractException, TimeoutException, InterruptedException {
		Contract contract = network.getContract(contractName);
		byte[] result = contract.createTransaction("AddNewUserItem")
		.setEndorsingPeers(network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
		.submit(newUser.pkUser, newUser.coinNum.toString(), newUser.credit.toString(), "consumer");
		return new String(result, StandardCharsets.UTF_8);
	}
	
	public static Integer QueryCoinNum(String pkUser) throws ContractException {
		Contract contract = network.getContract(contractName);
		byte[] result = contract.evaluateTransaction("QueryCoinNum", pkUser);
		return Integer.parseInt(new String(result, StandardCharsets.UTF_8));
	}
	
	public static Float QueryCredit(String pkUser) throws ContractException {
		Contract contract = network.getContract(contractName);
		byte[] result = contract.evaluateTransaction("QueryCredit", pkUser);
		return Float.parseFloat(new String(result, StandardCharsets.UTF_8));
	}
	
	public static String ChangeRole(String pkUser, String role) throws ContractException, TimeoutException, InterruptedException {
		Contract contract = network.getContract(contractName);
		byte[] result = contract.createTransaction("AddNewUserItem")
		.setEndorsingPeers(network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
		.submit(pkUser, "provider");
		return new String(result, StandardCharsets.UTF_8);
	}
	
}