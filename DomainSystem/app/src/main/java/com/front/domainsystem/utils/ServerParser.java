package com.front.domainsystem.utils;

import android.util.Log;

import org.json.JSONObject;

public class ServerParser {

    // JSON Node names
    private static final String TAG_IP = "servers";
    private static final String TAG_SSL_GRADE = "ssl_grade";
    private static final String TAG_COUNTRY = "country";
    private static final String TAG_OWNER = "owner";

    // Server attributes
    public String ipAddress;
    public String sslGrade;
    public String country;
    public String owner;

    // Constructor which allow us to map a JSON to the server attributes
    public ServerParser(JSONObject jsonObject){
        try{
            ipAddress=jsonObject.getString(TAG_IP);
            sslGrade= jsonObject.getString(TAG_SSL_GRADE);
            country=jsonObject.getString(TAG_COUNTRY);
            owner=jsonObject.getString(TAG_OWNER);
        }catch (Exception e){
            Log.e("Domain Parser", "Error parsing data " + e.toString());
        }
    }

    //Creates the string to show
    @Override
    public String toString() {
        String result="";
        result+="Servidor que tiene una dirección ip: "+ipAddress;

        if(!sslGrade.equals("")){
            result+="Su grado SSL calificado por SSLabs es: "+sslGrade+"\n";
        }

        if(!country.equals("")){
            result+="El servidor se encuentra en: "+country+"\n";
        }

        if(!owner.equals("")){
            result+="El dueño de la IP es: "+owner+"\n";
        }

        return result;
    }

}
