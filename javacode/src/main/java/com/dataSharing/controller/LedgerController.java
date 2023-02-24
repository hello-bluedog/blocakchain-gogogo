package com.dataSharing.controller;

import org.hyperledger.fabric.gateway.ContractException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.dataSharing.entity.ledgerEntity.CoinNumAndCredit;
import com.dataSharing.service.FabricService;

import io.grpc.netty.shaded.io.netty.handler.timeout.TimeoutException;

@RestController
public class LedgerController {
	@Autowired 
	public FabricService fabricService;
	
	
	@RequestMapping("/add")
	public String AddNewUserItem(String pkUser, int coinNum, float credit) throws ContractException, TimeoutException, InterruptedException, java.util.concurrent.TimeoutException {
		CoinNumAndCredit newUser = new CoinNumAndCredit(pkUser, coinNum, credit);
		return fabricService.AddNewUserItem(newUser);
	}
	
	@GetMapping("/query")
	public String QueryCoinNum(String pkUser) throws ContractException {
		return fabricService.QueryCoinNum(pkUser);
	}
	@GetMapping("/gg")
	public String hhh() {
		return "hello";
	}
	
}
