package com.dataSharing.entity.ledgerEntity;

import lombok.Data;

@Data
public class CoinNumAndCredit {
	public String pkUser;
	public Integer coinNum;
	public Float credit;
	public String role;
	public CoinNumAndCredit(String pkUser, Integer coinNum, Float credit, String role) {
		super();
		this.pkUser = pkUser;
		this.coinNum = coinNum;
		this.credit = credit;
		this.role = role;
	}
}