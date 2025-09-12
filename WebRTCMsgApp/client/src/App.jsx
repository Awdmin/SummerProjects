import { useState, useRef } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";

function App() {
  const peerRef = useRef(null);
  const targetRef = useRef(null);
  const dataChannelRef = useRef(null);

  const [users, setUsers] = useState([]);
  const [user, setUser] = useState("");
  const [registered, setRegistered] = useState(false);
  const [msgs, setMsgs] = useState([]);
  const [title, setTitle] = useState("WebRTC Chat");

  const socketURI = import.meta.env.REACT_APP_BACKEND_URI;
  const socket = new WebSocket(socketURI);

  socket.onconnectionstatechange = () => {
    console.log("Connected to Go WebSocket server");
  };

  socket.onmessage = async (event) => {
    try {
      const msg = JSON.parse(event.data);
      switch (msg.type) {
        case "user-list":
          setUsers(msg.data);
          break;

        case "offer":
          handleOffer(msg);
          console.log("Received offer");
          break;

        case "answer":
          await peerRef.current.setRemoteDescription(
            new RTCSessionDescription(msg.data),
          );
          console.log("Received answer");
          break;

        case "candidate":
          if (msg.data) {
            await peerRef.current.addIceCandidate(
              new RTCIceCandidate(msg.data),
            );
            console.log("Received candidate");
          }
          break;

        case "error":
          if (msg.data === "UAE") {
            setRegistered(false);
          }
          break;

        default:
          console.log("Unknown message type:", msg.type);
      }
    } catch (error) {
      console.error("Error parsing message:", error);
    }
  };

  const registerUser = () => {
    const name = document.getElementById("Name").value;
    const msg = {
      type: "register",
      from: name,
    };
    socket.send(JSON.stringify(msg));
    setRegistered(true);
    setUser(name);
  };

  const handleUserClick = async (targetUserId) => {
    targetRef.current = targetUserId;
    setTitle(`Chat with ${targetUserId}`);
    peerRef.current = createPeerConnection();

    const dataChannel = peerRef.current.createDataChannel("chat");
    setupDataChannel(dataChannel);

    const offer = await peerRef.current.createOffer();
    await peerRef.current.setLocalDescription(offer);

    socket.send(
      JSON.stringify({
        type: "offer",
        from: user,
        to: targetUserId,
        data: offer,
      }),
    );
  };

  const handleOffer = async (msg) => {
    console.log(msg);

    targetRef.current = msg.from;
    setTitle(`Chat with ${targetRef.current}`);
    peerRef.current = createPeerConnection();

    await peerRef.current.setRemoteDescription(
      new RTCSessionDescription(msg.data),
    );

    const answer = await peerRef.current.createAnswer();
    await peerRef.current.setLocalDescription(answer);

    socket.send(
      JSON.stringify({
        type: "answer",
        from: user,
        to: targetRef.current,
        data: answer,
      }),
    );
  };

  const createPeerConnection = () => {
    const pc = new RTCPeerConnection();
    pc.onicecandidate = (event) => {
      if (event.candidate) {
        socket.send(
          JSON.stringify({
            type: "candidate",
            from: user,
            to: targetRef.current,
            data: event.candidate,
          }),
        );
      }
    };

    pc.onconnectionstatechange = () => {
      console.log("Connection state:", pc.connectionState);
    };

    pc.ondatachannel = (event) => {
      const dataChannel = event.channel;
      setupDataChannel(dataChannel);
    };

    return pc;
  };

  const setupDataChannel = (dataChannel) => {
    dataChannelRef.current = dataChannel;

    dataChannel.onmessage = (event) => {
      console.log("Received message:", event.data);
      setMsgs((prev) => [
        ...prev,
        { user: targetRef.current, to: user, msg: event.data },
      ]);
    };

    dataChannel.onopen = () => {
      console.log("Data channel is open");
    };

    dataChannel.onclose = () => {
      console.log("Data channel is closed");
    };
  };

  const handleSend = () => {
    const data = document.getElementById("Msg").value;
    if (dataChannelRef.current) {
      dataChannelRef.current.send(data);
    }
    setMsgs((prev) => [
      ...prev,
      { user: user, to: targetRef.current, msg: data },
    ]);
    document.getElementById("Msg").value = "";
  };

  return (
    <>
      {!registered && (
        <div className="register">
          register:
          <input type="text" placeholder="Name" id="Name" />
          <button onClick={registerUser}>register</button>
        </div>
      )}
      <div className="container">
        <div className="usrContainer">
          {users
            .filter((usr) => usr !== user)
            .map((usr) => (
              <div
                className="user"
                key={usr}
                onClick={() => handleUserClick(usr)}
              >
                {usr}
              </div>
            ))}
        </div>
        <div className="chatContainer">
          <div className="title">{title}</div>
          <div className="msgArea">
            {" "}
            {msgs
              .filter(
                (msg) =>
                  msg.user === targetRef.current ||
                  (msg.user === user && msg.to === targetRef.current),
              )
              .map((msg, index) => (
                <div
                  className="msg"
                  key={index}
                  style={{
                    backgroundColor: msg.user === user ? "#999999" : "#17153b",
                  }}
                >
                  {msg.user}: {msg.msg}
                </div>
              ))}
          </div>
          <div className="inputArea">
            <input type="text" placeholder="Msg" id="Msg" />
            <button onClick={handleSend}>send</button>
          </div>
        </div>
      </div>
    </>
  );
}

export default App;
