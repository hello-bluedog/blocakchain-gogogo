package com.dataSharing.mapper;

import java.util.List;

import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Select;
import org.apache.ibatis.annotations.Update;

import com.dataSharing.entity.dataSharingEntity.SkRequest;

@Mapper
public interface SkRequestMapper {
	@Insert("insert into SkRequest values(null, #{pkOfConsumer}, #{pkOfProvider}, null, 0)")
	public int newSkRequest(String pkOfConsumer, String pkOfProvider);
	
	@Select("select * from SkRequest where pkOfconsumer = #{pkOfConsumer}")
	public List<SkRequest> getRequest(String pkOfConsumer); 
	
	@Update("update SkRequest set cipherOfSk = #{cipherOfSk} where pkOfconsumer = #{pkOfConsumer} and pkOfProvider = #{pkOfProvider}")
	public int uploadCipherSk(String pkOfConsumer, String pkOfProvider, String cipherOfSk);
	
	@Select("select cipherOfSk from SkRequest where where pkOfconsumer = #{pkOfConsumer} and pkOfProvider = #{pkOfProvider} ")
	public String getCipherSk(String pkOfConsumer, String pkOfProvider);
}
