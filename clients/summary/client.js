const baseUrlSummary="https://api.sumsumsummary.me/v1/summary?url=https://www.google.com";
const baseUrl="https://api.sumsumsummary.me";
const registerRoute="/v1/users";
const loginRoute="/v1/sessions";
var ws;
var authToken="abcd";
var thisUser={};
let headerAuthorization="Authorization"
let printlnResponse=(responseBody)=>{
  console.log("Println Response");
  console.log(responseBody);
}
let getUserIdFromInput=()=>{
  return parseInt(document.getElementById("userIdInput").value);
}
let makeRandomChannelJson=()=>{
  return {
    id:Math.floor(Math.random()*30000),
    name:getUserNameValue(),
    description:"This is a random discription",
    private:false,
    members: [getUserIdFromInput()]
  };
}
var channelRoutes=[{
  route:"/v1/channels",
  requestType:"POST",
  body:{},
  preflightSetBody:makeRandomChannelJson,
  callback:printlnResponse
},{
  route:"/v1/channels/",
  routeSuffix:true,
  requestType:"POST",
  preflightSetBody:makeRandomChannelJson,
  callback:printlnResponse
},{
  route:"/v1/channels/",
  routeSuffix:true,
  requestType:"PATCH",
  body:{
    id:44,
    name:"newName44",
    description:"newDescription",
    private:false,
    members:[1,2,3,thisUser.id]
  },
  callback:printlnResponse
},{
  route:"/v1/channels/",
  routeSuffix:true,
  requestType:"DELETE",
  body:{},
  preflightSetBody:makeRandomChannelJson,
  callback:printlnResponse
},{
  route:"/v1/channels/",
  memberRouteSuffix:true,
  requestType:"POST",
  body:{
    //should be a new members
  },
  callback:printlnResponse
},{
  route:"/v1/channels/",
  memberRouteSuffix:true,
  requestType:"DELETE",
  body:{
    //should be the deleted members identifiers
  },
  callback:printlnResponse
},{
  route:"/v1/messages/",
  messageSuffix:true,
  requestType:"PATCH",
  body:{
    //should be a new members
  },
  callback:printlnResponse
},{
  route:"/v1/messages/",
  messageSuffix:true,
  requestType:"DELETE",
  body:{
    //should be a new members
  },
  callback:printlnResponse
}];
let getChannelNameInputValue=()=>{
  let val=document.getElementById("channelNameInput").value;
  if(val==""){
    val="noID";
  }
  return val;
}
let credentials= {
	email:"abcd",
	password:"password"
}

let onLoginResponse=(responseBody)=>{
  console.log("Login Response");
  thisUser=responseBody;
  console.log(responseBody);
}
function populateRouteOptionsSelect(){
  let postRouteSelect=document.getElementById("postRouteSelect");
  for(let i=0;i<channelRoutes.length;i++){
    let newOption=document.createElement("option");
    let optionObj=channelRoutes[i];
    newOption.innerHTML=optionObj.requestType+" "+optionObj.route;
    if(optionObj.routeSuffix){
      newOption.innerHTML=newOption.innerHTML+"{channelId}";
    }else if(optionObj.memberRouteSuffix){
      newOption.innerHTML=newOption.innerHTML+"{channelId}/members";
    }else if(optionObj.messageSuffix){
      newOption.innerHTML=newOption.innerHTML+"{messageId}";
    }
    postRouteSelect.appendChild(newOption);
  }
}
let attemptLogin=()=>{
  //Attempts a login fetch request.
  let theseCredentials={
    email:document.getElementById("loginInput").value,
    password:"password"
  }
  makePostCall(baseUrl+loginRoute,onLoginResponse,JSON.stringify(theseCredentials));
}

let makeNewUserRequest=()=>{
  let username=getUserNameValue();
  console.log("attempting login");
  let payload={
    email:username,
    password:"password",
    passwordConf:"password",
    userName:username,
    firstName:"firstName",
    lastName:"lastName"
  };
  urlStr=baseUrl+registerRoute;
  makePostCall(urlStr,onRegisterResponse,JSON.stringify(payload));
}
function routeResponseCallBack(responseBody){
  console.log("routeResponseCallBack");
  console.log(responseBody);
}
function makeRouteRequest(optionIndex){
  let optionObj=channelRoutes[optionIndex];
  let requestType=optionObj.requestType;
  urlStr=baseUrl+optionObj.route;
  if(optionObj.routeSuffix){
    urlStr=urlStr+getChannelNameInputValue();
  }else if(optionObj.memberRouteSuffix){
    urlStr=urlStr+getChannelNameInputValue()+"/members";
  }else if(optionObj.messageSuffix){
    urlStr=urlStr+getChannelNameInputValue();
  }
  if(requestType=="GET"){
    makeGetCall(urlStr,routeResponseCallBack);
  }else{
    let payload=optionObj.body;
    if(optionObj.preflightSetBody){
      //console.log("preforming custom setbody for post request");
      payload=optionObj.preflightSetBody();
    }
    payload.auth=authToken;
    console.log("Sending payload",payload);
    console.log("To "+urlStr)
    if(requestType=="POST"){
      makePostCall(urlStr,routeResponseCallBack,JSON.stringify(payload));
    }else if(requestType=="DELETE"){
      makeDeleteCall(urlStr,routeResponseCallBack,JSON.stringify(payload));
    }else if(requestType=="PATCH"){
      makePatchCall(urlStr,routeResponseCallBack,JSON.stringify(payload));
    }else{
      console.log("Invalid option type"+requestType);
    }
  }
}
function sendTestPayload(payloadStr){
  let payloadTest={
    message:document.getElementById("websocketPayloadInput").value
  }
  sendJSON(payloadTest);
}
function sendJSON(payloadJSON){
  ws.send(JSON.stringify(payloadJSON),()=>{
    console.log("Message Sent");
  });
}
function attemptConnectToWebsocket(authString){
    ws= new WebSocket('wss://api.sumsumsummary.me/ws?auth='+authString);
    console.log("Attempting ping start");
    ws.onopen = function () {
      ws.send(JSON.stringify({message:"Ping"}),()=>{
      console.log("ping end");
    });
    };
}
function setBaseOnMessage(){
  ws.onmessage = function (e) {
    console.log("From Server:"+ e.data);
  };
}
let onRegisterResponse=function(responseBody){
  console.log("New user response"+responseBody);
  thisUser=responseBody;
  document.getElementById("titleOutput").innerHTML="Register Response";
  document.getElementById("mainOutput").innerHTML=responseBody.firstName+" "+responseBody.lastName;
}

//
let summaryGetResponse=function(responseJSON){
  responseBody=responseObj.body;
  console.log("response"+responseBody.description);
  document.getElementById("titleOutput").innerHTML=responseBody.title;
  document.getElementById("mainOutput").innerHTML=responseBody.description;
}

/* Preforms a fetch call on the server and calls the given function.
@param {string} urlSuffix - the string to append to the base api link.
@param {function} callFunction - the function to call on a valid request.
@param {object} params - The request's parameters. */
function doFetch(URL, callFunction, params) {
  params.mode = "cors";

  //Be prepared to comment this out, causes CORS errors
  /*params.headers = new Headers({
    "Access-Control-Allow-Origin": "*",
    "Access-Control-Allow-Headers": "*",
    "Content-Type":"application/json"
  });*/
  /*params.header = {//This works but doesn't send the content-type
    "Access-Control-Allow-Origin": "*",
    "Access-Control-Allow-Headers": "*",
    "Content-Type":"application/json"
  };*/
  params.headers = {//Causes CORS error
    "Access-Control-Allow-Origin": "*",
    "Access-Control-Allow-Headers": "*",
    "Access-Control-Expose-Headers":"*",
    "Content-Type":"application/json"//,
  //  'Authorization': 'Bearer ' + authToken
};
  fetch(URL, params)
    .then(checkStatus)
    .then(returnBody)
    .then(callFunction)
    .catch(errorMessage);
}
let returnBody=function(responseBody){
  console.log(responseBody);
  return responseBody;
}

/* Makes GET call to server and calls the given function using the returned data.
@param {string} urlSuffix - the string to append to the base api url.
@param {function} callFunction - the function to call on a properly executed request. */
function makeGetCall(url, callFunction) {
  doFetch(url, callFunction, {
    method: "GET"
  });
}

//BODY should be a json String
function makePostCall(url, callFunction,body) {
  doFetch(url, callFunction, {
    method: "POST",
    body:body
  });
}
//BODY should be a json String
function makeDeleteCall(url, callFunction,body) {
  doFetch(url, callFunction, {
    method: "DELETE",
    body:body
  });
}
//BODY should be a json String
function makePatchCall(url, callFunction,body) {
  doFetch(url, callFunction, {
    method: "PATCH",
    body:body
  });
}

  /*Tells the user that an error has occured with a fetch request through the console.
  @param {string} statusMessage - The message to display */
  function errorMessage(statusMessage) {
    console.log("Error "+statusMessage);
  }

  /* Checks the response status and if valid returns the response.
  If invalid prints an error message to the console.
  @param {object} response - Response from server
  @return {string} - The server's response. */
  function checkStatus(response) {
    if (response.status >= 200 && response.status < 300) {
      let authHeader = response.headers.get(headerAuthorization);
      if(authHeader){
        console.log("header auth token is "+authHeader);
        //authToken=authHeader.substring(authHeader.indexOf("Bearer ")+"Bearer ".length);
        authToken=authHeader;
        setAuthTokenDisplay();
      }else{
        console.log("no header auth token found:");
      }
      return response.json();
    } else {
      return Promise.reject(new Error(response.status + ":" + response.statusText));
    }
  }
  function setAuthTokenDisplay(){
    document.getElementById("AuthTokenDisplay").innerHTML=authToken;
  }
  function getUserNameValue(){
    return document.getElementById("loginInput").value;
  }
