const invisibleClass="invisible";
const attemptLoginRoute="/v1/sessions";
const makeNewUserRoute="/v1/users";
const getOpenChannelsRoute="/v2/openchannels";
var currentChatId=null;


function setMinimizeButtonOnclicks(minSection,parentSection){
  let allMinimizeButtons=parentSection.querySelectorAll(":scope .hideXBox");
  for(let i=0;i<allMinimizeButtons.length;i++){
    let minButton=allMinimizeButtons[i];
    minButton.onclick=()=>{
      if(minSection.classList.contains(invisibleClass)){
        minSection.classList.remove(invisibleClass);
      }else{
        minSection.classList.add(invisibleClass);
      }
    }
  }
}
function getEmailFromField(){
  return document.getElementById("emailInput").value;
}
function getPasswordFromField(){
  return document.getElementById("passwordInput").value;
}
function onNewUserCreatedResponse(jsonData){
  console.log("New user created and logged in",jsonData);
  onLogin(jsonData);
}
function setDebateMessageHandler(){
  ws.onmessage = function (e) {
    let jsonObj=JSON.parse(e.data);
    console.log("From Server:", jsonObj);
    if(jsonObj.messageType=="debateMessage"){
      appendNewMessage(jsonObj);
    }
  };
}
function appendNewMessage(messageObj){
  if(messageObj.channelId==currentChatId){
  let parentElement=document.getElementById("appendMessagesTo");
  let baseDiv=document.createElement("div");
  baseDiv.classList.add("messageDiv");
  let innerDiv=document.createElement("div");
  let innerSpan=document.createElement("span");
  innerSpan.innerHTML="User:";
  innerDiv.appendChild(innerSpan);
  let innerSpan2=document.createElement("span");
  innerSpan2.innerHTML=messageObj.handle;
  innerSpan.appendChild(innerSpan2);
  baseDiv.appendChild(innerDiv);
  let innerDiv2=document.createElement("div");
  baseDiv.appendChild(innerDiv2);
  let paragraph=document.createElement("p");
  paragraph.innerHTML=messageObj.message;
  innerDiv2.appendChild(paragraph);
  parentElement.appendChild(baseDiv);
}else{
  console.log('invalid chat id '+messageObj.channelID);
}
}
function onLogin(jsonData){
  attemptConnectToWebsocket(authToken);
  setDebateMessageHandler();
  document.getElementById("loginStatusDiv").classList.remove("invisible");
  document.getElementById("loginUsernameSpan").innerHTML="Logged in as "+jsonData.userName;
}
function makeNewUserButtonOnclick(){
  document.getElementById("makeNewUserButton").onclick=()=>{
    let email=getEmailFromField();
    let password=getPasswordFromField();
    let validEmailRegex = /^.+@.+.com$/;

    if(!email.match(validEmailRegex)){
      document.getElementById("loginUsernameSpan").innerHTML="Invalid email, needs ...@...com";
      console.log("Invalid email, needs ...@...com");
    }
    if(password.length<7){
      document.getElementById("loginUsernameSpan").innerHTML="Invalid password, needs at least 6 charachters";
      console.log("Email is too short!!!");
    }
    console.log("attempting login");
    let payload={
      email:email,
      password:password,
      passwordConf:password,
      userName:email,
      firstName:document.getElementById("firstNameInput").value,
      lastName:document.getElementById("lastNameInput").value
    };
    urlStr=baseUrl+makeNewUserRoute;
    console.log("sending",urlStr,payload);
    makePostCall(urlStr,onNewUserCreatedResponse,JSON.stringify(payload));
  }
}
function makeLoginButtonOnclick(){
  document.getElementById("attemptLoginButton").onclick=()=>{
    let email=getEmailFromField();
    let password=getPasswordFromField();
    console.log("attempting login");
    let payload={
      email:email,
      password:password,
      passwordConf:password
    };
    urlStr=baseUrl+attemptLoginRoute;
    makePostCall(urlStr,onLogin,JSON.stringify(payload));
  }
}
function getChatDataToPopulateMessages(chatId,name){
  //Sets the chatid as a var
  if(authToken&&authToken!="abcd"){
  currentChatId=chatId;
  document.getElementById("debateNameSpan").innerHTML=name;
  document.getElementById("appendMessagesTo").innerHTML="";
}else{
  document.getElementById("debateNameSpan").innerHTML="Must log in debate.";
}
}
function onGetOpenChannelsCallback(jsonObj){
  let debateSelector=document.getElementById("debateSelector");
  debateSelector.innerHTML="";
  for(let i=0;i<jsonObj.length;i++){
    let baseDiv=document.createElement("div");
    let innerButton=document.createElement("button");
    innerButton.onclick=(onclickEvent)=>{
      getChatDataToPopulateMessages(jsonObj[i].id,jsonObj[i].name);
    };
    innerButton.innerHTML=jsonObj[i].name;
    baseDiv.appendChild(innerButton);
    debateSelector.appendChild(baseDiv);
  }
}
function populateDefaultChannels(){
  onGetOpenChannelsCallback([{id:0,name:"General"},{id:1,name:"Fake meat: Yay or Nay?"}]);
}
function populateChannelsDiv(){
  let debateSelector=document.getElementById("debateSelector");
  populateDefaultChannels();
  let urlStr=baseUrl+getOpenChannelsRoute;
  //makeGetCall(urlStr,onGetOpenChannelsCallback);
}
function setSendButtonOnclick(){
  document.getElementById("sendButton").onclick=()=>{
    if(currentChatId){
      let payload={
        messageType:"debateMessage",
        channelId:currentChatId,
        message:document.getElementById("messageInput").value+"",
        handle:document.getElementById("handleInput").value+"",
        username:thisUser.id+""
      };
      sendJSON(payload);
    }else{
      console.log("Current chat id is invalid.")
    }
  }
}
window.onload=()=>{
  setMinimizeButtonOnclicks(document.querySelector(".verticalInputDiv"),document.querySelector(".loginSection"));
  makeNewUserButtonOnclick();
  makeLoginButtonOnclick();
  populateChannelsDiv();
  setSendButtonOnclick();
}
