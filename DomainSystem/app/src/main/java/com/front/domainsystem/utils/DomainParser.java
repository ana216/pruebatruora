package com.front.domainsystem.utils;

import android.util.Log;

import org.json.JSONArray;
import org.json.JSONObject;

import java.util.ArrayList;

public class DomainParser{

    // JSON Node names
    private static final String TAG_SERVERS = "servers";
    private static final String TAG_SERVERS_CHANGED = "servers_changed";
    private static final String TAG_SSL_GRADE = "ssl_grade";
    private static final String TAG_PREVIOUS_SSLGRADE = "previous_ssl_grade";
    private static final String TAG_LOGO = "logo";
    private static final String TAG_TITLE = "title";
    private static final String TAG_IS_DOWN = "is_down";

    // Domain attributes
    public ArrayList<ServerParser> servers;
    public boolean serversChanged;
    public String sslGrade;
    public String previous_ssl_grade;
    public String logo;
    public String title;
    public boolean isDown;

    // Constructor which allow us to map a JSON to the domain attributes
    public DomainParser(JSONObject jsonObject) {
        servers= new ArrayList<ServerParser>();
        try {
            // Getting Array of Servers
            JSONArray tmpServerArray = jsonObject.getJSONArray(TAG_SERVERS);
            // looping through All Servers
            for(int i = 0; i < tmpServerArray.length(); i++){
                JSONObject JSONserver = tmpServerArray.getJSONObject(i);
                ServerParser server= new ServerParser(JSONserver);
                servers.add(server);
            }
            //Getting the remaining domain attributes
            serversChanged=jsonObject.getBoolean(TAG_SERVERS_CHANGED);
            sslGrade=jsonObject.getString(TAG_SSL_GRADE);
            previous_ssl_grade=jsonObject.getString(TAG_PREVIOUS_SSLGRADE);
            logo= jsonObject.getString(TAG_LOGO);
            title=jsonObject.getString(TAG_TITLE);
            isDown= jsonObject.getBoolean(TAG_IS_DOWN);

        }catch (Exception e){
            Log.e("Domain Parser", "Error parsing data " + e.toString());
        }

    }

    //Creates the string to show
    @Override
    public String toString() {
        String result="";
        if(this.isDown){
            result+="El dominio está caído o no existe";
        }else{
            if(this.serversChanged){
                result+="Ha cambiado de servidores hace menos de una hora.\n";
            }else{
                result+="Sus servidores no han cambiado en la última hora.\n";
            }
            if(!this.sslGrade.equals("")){
                result+="El grado más bajo de todos los servidores es: "+this.sslGrade+"\n";
            }
            if(!this.previous_ssl_grade.equals("")){
                result+="El grado más bajo de todos los servidores es: "+this.previous_ssl_grade+"\n";
            }
            if(!this.logo.equals("")){
                result+="La URL del logo usado en la página es: "+this.logo+"\n";
            }
            if(!this.title.equals("")){
                result+="El título de la página es: "+this.title+"\n";
            }
            int numServers= servers.size();
            if(numServers!=0){
                result+="El dominio no tiene servidores";
            }else{
                result+="El dominio tiene "+numServers+" servidor(es) \n";
                for (int i=0;i<numServers;i++){
                    result+="Servidor "+i+"\n";
                    result+=servers.get(i).toString();

                }
            }

        }
        return result;
    }
}
