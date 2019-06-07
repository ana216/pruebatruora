package com.front.domainsystem;

import org.junit.Test;

import static org.junit.Assert.assertTrue;


//Class which test the main functionalities of the system: show domain info and show history
public class MainTest {

    private MainActivity scenary;

    public MainTest(){
        scenary = new MainActivity();
    }

    @Test
    public void useShowInfo() {
        String test= scenary.showInfoDomain("truora.com");
        assertTrue(test!=null&&!test.equals(""));

    }

    @Test
    public void useShowHistory() {
        String test= scenary.showHistoryDomains("truora.com");
        assertTrue(test!=null&&!test.equals(""));

    }
}
