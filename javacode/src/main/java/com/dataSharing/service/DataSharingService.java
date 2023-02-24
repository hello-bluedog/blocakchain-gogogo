package com.dataSharing.service;

import java.util.concurrent.TimeoutException;

import org.hyperledger.fabric.gateway.ContractException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.dataSharing.mapper.SkRequestMapper;

@Service
public class DataSharingService {
	public static float creditThreshold = (float)2.0;
	public static int coinNumThreshold = 2;
	@Autowired
	public SkRequestMapper skRequestMapper;
	public boolean beProvider(String pkUser) throws ContractException, TimeoutException, InterruptedException {
		float credit = FabricService.QueryCredit(pkUser);
		Integer coinNum = FabricService.QueryCoinNum(pkUser);
		if(credit < creditThreshold && coinNum < coinNumThreshold) {
			return false;
		}
		//change role in ledger
		FabricService.ChangeRole(pkUser, "provider");
		return true;
	}
	
	public boolean dataRequest(String pkOfConsumer, String pkOfProvider) {
		if(skRequestMapper.newSkRequest(pkOfConsumer, pkOfProvider) == 1) {
			return true;
		}
		return false;
	}
	
	public int skSend(String pkOfConsumer, String skOfData) {
		return skRequestMapper.uploadCipherSk(pkOfConsumer, pkOfConsumer, skOfData);
	}
	
	public String skGet(String pkOfConsumer) {
		return skRequestMapper.getCipherSk(pkOfConsumer, pkOfConsumer);
	}
}
