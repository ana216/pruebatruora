package com.front.domainsystem;

import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.content.ContextCompat;

import android.Manifest;
import android.content.pm.PackageManager;
import android.os.Bundle;

import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import com.front.domainsystem.utils.DomainParser;
import com.front.domainsystem.utils.HttpGetRequest;
import com.front.domainsystem.utils.ServerParser;


import org.json.JSONArray;
import org.json.JSONObject;

public class MainActivity extends AppCompatActivity {

    //Endpoints URL
    public static final String ENDPOINT_INFO_DOMAIN="http://10.0.2.2:2020/servers/";
    public static final String ENDPOINT_DOMAIN_HISTORY="http://10.0.2.2:2020/servers/alldomains";

    //Permission codes
    private static final int MY_INTERNET_PERMISSION_CODE = 200;

    //JSON Node names
    private static final String TAG_ITEMS= "items";

    // UI components
    private Button btn_search_history;
    private Button btn_search_domainInfo;
    private EditText etxt_domain_info;
    private TextView txt_domain_history;
    private TextView txt_servers_domain_info;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        btn_search_history= findViewById(R.id.btn_search_history);
        btn_search_domainInfo= findViewById(R.id.btn_select_all_server);
        etxt_domain_info=findViewById(R.id.text_domainName_info);
        txt_domain_history= findViewById(R.id.text_review_history);
        txt_servers_domain_info= findViewById(R.id.text_server_domain);

        btn_search_domainInfo.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                try{
                    String domainToFind= etxt_domain_info.getText().toString().replace(" ","");
                    if(!domainToFind.equals("")){
                        String url=ENDPOINT_INFO_DOMAIN+domainToFind;
                        txt_servers_domain_info.setText(showInfoDomain(url));

                    }else{
                        Toast toast = Toast.makeText(getApplicationContext(), "Debe ingresar un dominio válido", Toast.LENGTH_SHORT);
                        toast.show();
                    }

                }catch (Exception e){
                    Toast toast = Toast.makeText(getApplicationContext(), e.getMessage(), Toast.LENGTH_SHORT);
                    toast.show();
                    Log.e("Main", e.getMessage());
                }
            }
        });


        btn_search_history.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {

                txt_domain_history.setText(showHistoryDomains(ENDPOINT_DOMAIN_HISTORY));
            }
        });

    }

    //Method which returns a String corresponding to the info parsed of the json got from the endpoint that allow to get the info of each domain
    public String showInfoDomain( String Url){
        String result="";
        if (ContextCompat.checkSelfPermission(getApplicationContext(), Manifest.permission.INTERNET)
                != PackageManager.PERMISSION_GRANTED) {
            requestPermissions(new String[]{Manifest.permission.INTERNET},
                    MY_INTERNET_PERMISSION_CODE);
        } else {
            try{
                HttpGetRequest getRequest=new HttpGetRequest();
                String resGet=getRequest.execute(Url).get();
                JSONObject jsonObject= new JSONObject(resGet);
                DomainParser domainParser=new DomainParser(jsonObject);
                result="Información del dominio:\n"+domainParser.toString();
            }catch(Exception e){
                Toast toast = Toast.makeText(getApplicationContext(), e.getMessage(), Toast.LENGTH_SHORT);
                toast.show();
                Log.e("Main", e.getMessage());

            }
        }
        return result;
    }

    //Method which returns a String corresponding to the info parsed of the json got from the endpoint that allows to get the info about recently domains consulted
    public String showHistoryDomains(String url) {
        String result = "";
        if (ContextCompat.checkSelfPermission(getApplicationContext(), Manifest.permission.INTERNET)
                != PackageManager.PERMISSION_GRANTED) {
            requestPermissions(new String[]{Manifest.permission.INTERNET},
                    MY_INTERNET_PERMISSION_CODE);
        } else {
            try {
                HttpGetRequest getRequest = new HttpGetRequest();
                String resGet = getRequest.execute(url).get();
                JSONObject jsonObject= new JSONObject(resGet);
                JSONArray tmpDomainArray =jsonObject.getJSONArray(TAG_ITEMS);
                for(int i = 0; i < tmpDomainArray.length(); i++){
                    String[] infoDomain = tmpDomainArray.getString(i).split(" ");
                    result+="Dominio "+(i+1)+": "+infoDomain[0]+"\n";
                }
            } catch (Exception e) {
                Toast toast = Toast.makeText(getApplicationContext(), e.getMessage(), Toast.LENGTH_SHORT);
                toast.show();
                Log.e("Main", e.getMessage());
            }

        }
        return result;
    }

    //Method which handles permission issues
        @Override
        public void onRequestPermissionsResult ( int requestCode, @NonNull String[] permissions,
        @NonNull int[] grantResults){
            super.onRequestPermissionsResult(requestCode, permissions, grantResults);
            switch (requestCode) {
                case MY_INTERNET_PERMISSION_CODE:
                    if (grantResults[0] == PackageManager.PERMISSION_GRANTED) {
                        Toast.makeText(getApplicationContext(), "Permiso para acceder a internet aceptado", Toast.LENGTH_LONG).show();
                    } else {
                        Toast.makeText(getApplicationContext(), "Permiso para acceder a Internet denegado", Toast.LENGTH_LONG).show();
                    }

            }
        }

    }

