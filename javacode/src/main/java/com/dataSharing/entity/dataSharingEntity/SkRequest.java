package com.dataSharing.entity.dataSharingEntity;

public class SkRequest {
	public int id;
	public String pkOfConsumer;
	public String pkOfProvider;
	public String cipherOfSk;
	public String status;
	public SkRequest(int id, String pkOfConsumer, String pkOfProvider, String cipherOfSk, String status) {
		super();
		this.id = id;
		this.pkOfConsumer = pkOfConsumer;
		this.pkOfProvider = pkOfProvider;
		this.cipherOfSk = cipherOfSk;
		this.status = status;
	}
	
}
