# Master Debater
By Boya Cao, Brandon Chong, Yvonne Liu, Leon Tan
## Project Description 
We plan to build an online debating web application to let students, developers, and people of all software backgrounds to debate with one another on the internet. A new debate can be created by any registered user and any logged-in user will be able to join a debate room. Debates will last for a total of an hour. In the debate room, there will be a separation of roles: debaters, audience, and moderator(s); debaters will debate in a team format, audience will be able to read the discussions and vote, and the moderator(s) will help make sure there are no "bad egg" participants. Once the debate is over, the audience will vote affirmative or negative for on the debate claim. Concluded debates will be archived showing the winner and a log of the discussion.

Our target audiences including people of any software background and experience that wish to develop and test their knowledge across a wide span of topics, in debate form. For this audience, Master debater will be a great way to meet more people, both like-minded and opposite-minded, on the internet. By debating, we hope to help current and future software developers to expand their thinking horizon, so that they can be more informed in their development decisions.

Additionally, arguing with people in real life can be difficult because it is often difficult confronting others about their opinions. Also, a person's reputation can be damaged or a good relationship can be ruined due to a debate. But with Master Debater users can choose to debate anonymously (hidden behind a username) that would not affect them (unless their username was identifiable). By using Master Debater, we hope that the skills of debating, questioning, and collaborative brainstorming will transfer to real life and make the users better team players and developers.

From a utilitarian perspective the joy that these people generate from each other is a moral good and enabling them to fight with each other is just. We want to make the world a better place by educating others on various different CS and programming topics.
 
## Technical Description
### Infrastructure
The system will heavily rely on Docker containers to isolate the client, server, and database. Users will only interact with the client, and the server will perform all the required actions to provide a seamless user experience for debating.

![Master Debater architectural diagram mapping](./images/INFO441_proposal_diagram.png)

https://www.lucidchart.com/invitations/accept/c719ef91-0be0-462c-9432-0dbcb8a6383c

### User Stories
|Priority|User|Description|
|--------|----|-----------|
|P0|As a debater|I want to filter out chatrooms with my favoured topics. I want to be able to post within chat rooms. I want to see my mailbox for new messages. I want to be able to delete my account when I'm done with this app.<br />API routes: <ul><li>GET /chatroom/{query}</li><li>POST /chatroom/{id}/post</li><li>DELETE /user</li></ul>|
|P1|As an audience|I want to vote affirmative or negative for a claim in a chatroom. I want to be able to create an account.<br />API routes:<ul><li>POST /chatroom/{id}?{vote}</li><li>POST /user/create</li></ul>|
|P2|As a host|I want to create a chatroom with my favoured topics.<br />API routes:<ul><li>GET /chatroom/create</li><li>API:POST /chatroom/{id}</li></ul>|
|P3|As a moderator|I want to be able to kick out bad participants.<br />API routes:<ul><li>POST /chatroom/{id}/kick/{userId}</li></ul>|
 

### API Design

- **/api/user**: user control - GET current user information, POST new current user information, and DELETE current user.
    - GET: Get current user information
        - “200”: “application/json” - successfully got current user’s information and returns user information in JSON format.
        - “401”: “error” - cannot verify the current user
        - “500”: “error” - internal server error
    - POST: Update current user information
        - “201”: “application/json” - successfully updated current user’s information
        - “402”: “error” - cannot read body or incorrect body
        - “500”: “error” - internal server error
    - DELETE: Delete current user
        - “201”: “application/json” - successfully deleted current user’s information
        - “401”: “error” - cannot verify the current user
        - “500”: “error” - internal server error
- **/api/user/create**: POST new user
    - POST: Create new user
        - “200”: “application/json” - successfully created a new user, logs them in, and returns their information in JSON format.
        - “402”: “error” - “error” - cannot read body or incorrect body
        - “500”: “error” - internal server error
- **/api/user/login**: POST existing user credentials
    - POST: Login user
        - “200”: “application/json” - successfully logs user in, and returns their information in JSON format.
        - “401”: “error” - cannot verify the user’s credentials
        - “500”: “error” - internal server error
- **/api/chatroom**: GET a list of all current debates
    - GET: Get all current debates
        - “200”: “application/json” - successfully got all current debates and returns them in JSON format with name of debate, description, host, and persons count.
        - “500”: “error” - internal server error
- **/api/chatroom/create**: POST create a debate
    - POST: create a debate
        - “200”: “application/json” - successfully created and joined the debate (as either moderator, debater, or audience; this is the creator’s choice)
        - “401”: “error” - cannot verify the current user
        - “500”: “error” - internal server error
- **/api/chatroom/{id}**: POST join a debate
    - POST: Join a debate
        - “200”: “application/json” - successfully joined a debate (either as a debater or audience)
        - “401”: “error” - cannot verify the current user
        - “404”: “error” - could not find the debate
        - “500”: “error” - internal server error
- **/api/chatroom/{id}?{vote}**: POST vote on debate
    - POST: Vote on debate
        - “200”: “application/json” - successfully voted on debate
        - “401”: “error” - cannot verify the current user
        - “402”: “error” - cannot read body or incorrect body
        - “404”: “error” - could not find the debate
        - “500”: “error” - internal server error
- **/api/chatroom/{id}/kick/{userId}**: POST kick a participant
    - POST: Kick a participant
        - “200”: “application/json” - successfully kicked participant
        - “401”: “error” - cannot verify the current user
        - “404”: “error” - could not find the participant
        - “500”: “error” - internal server error
- **/api/chatroom/{id}/post**: POST a message
    - POST: Post a message
        - “200”: “application/json” - successfully posted a message
        - “401”: “error” - cannot verify the current user
        - “405”: “error” - could not post message due to not being the team’s turn
        - “500”: “error” - internal server error
- **/api/chatroom/archived**: GET a list of all archived debates
    - GET: Get all archived debates
        - “200”: “application/json” - successfully got all current debates
        - “500”: “error” - internal server error
- **/api/chatroom/archived/{id}**: GET a specific archived debate
    - GET: Get all archived debates
        - “200”: “application/json” - successfully got all archived debates
        - “404”: “error” - could not find the debate
        - “500”: “error” - internal server error
