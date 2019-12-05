
  window.onload=function(){
    console.log("sanity check:9999rrf9ws");
    setAuthTokenDisplay();
    document.getElementById("loginInput").value=Math.random()*2000+"@abc";
    document.getElementById("submitButton").onclick=function(){
      let urlVal=document.getElementById("urlInput").value;
      console.log("Button Click:"+urlVal);
      makeGetCall(urlVal,summaryGetResponse);
    };
    document.getElementById("makeNewUser").onclick=function(){
      console.log("New user button clicked:");
      makeNewUserRequest();
    }
    document.getElementById("submitLoginInput").onclick=function(){
      console.log("Attempt loginbutton clicked:");
      attemptLogin();
    }
    document.getElementById("connectButton").onclick=function(){
      console.log("Connect button clicked");
      attemptConnectToWebsocket(authToken);
      setBaseOnMessage();
    }
    document.getElementById("sendPayloadButton").onclick=function(){
      console.log("Send websocket payload");
      let websocketPayloadInput=document.getElementById("websocketPayloadInput");
      sendTestPayload(websocketPayloadInput.value);
    }
    populateRouteOptionsSelect();
    document.getElementById("submitFetchSelectorQuery").onclick=function(){
      console.log("Send fetch payload");
      let optionsSelector=document.getElementById("postRouteSelect");
      makeRouteRequest(optionsSelector.selectedIndex);
    }
  }
