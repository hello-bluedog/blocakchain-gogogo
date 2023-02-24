package com.dataSharing.entity.ledgerEntity;

import lombok.Data;

@Data
public class CipherIndex {
	public String pkEnc;
	public String pkUser;
	public String roadName;
	public Integer messageLevel;
	public boolean weather;
	public boolean condition;
	public boolean traffic;
	public boolean averageSpeed;
	public CipherIndex(String pkEnc, String pkUser, String roadName, Integer messageLevel, boolean weather,
			boolean condition, boolean traffic, boolean averageSpeed) {
		super();
		this.pkEnc = pkEnc;
		this.pkUser = pkUser;
		this.roadName = roadName;
		this.messageLevel = messageLevel;
		this.weather = weather;
		this.condition = condition;
		this.traffic = traffic;
		this.averageSpeed = averageSpeed;
	}

}