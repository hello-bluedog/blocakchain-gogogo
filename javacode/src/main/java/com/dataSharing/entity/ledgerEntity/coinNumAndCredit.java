package com.dataSharing.entity.ledgerEntity;

import lombok.Data;

@Data
public class coinNumAndCredit {
	public String pkUser;
	public Integer coinNum;
	public Float credit;
	public coinNumAndCredit(String pkUser, Integer coinNum, Float credit) {
		super();
		this.pkUser = pkUser;
		this.coinNum = coinNum;
		this.credit = credit;
	}
}
