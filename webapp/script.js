document.addEventListener("DOMContentLoaded", () => {
  const betForm = document.getElementById("bet-form");
  const submitButton = document.getElementById("sumbit-bet");
  const getBlockchainButton = document.getElementById("get-blockchain");
  const blockchainData = document.getElementById("blockchain-data");
  const mineBlockButton = document.getElementById("mine-block");
  const mineResult = document.getElementById("mine-result");

  submitButton.addEventListener("click", (event) => {
    //event.preventDefault();

    const playerName = betForm.getElementById("playername").value;
    const matchID = betForm.getElementById("matchid").value;
    const teamOneScore = betForm.getElementById("teamonescore").value;
    const teamTwoScore = betForm.getElementById("teamtwoscore").value;
    console.log(playerName);

    const bet = {
      playername: playerName,
      matchid: matchID,
      teamonescore: parseInt(teamOneScore, 10),
      teamtwoscore: parseInt(teamTwoScore, 10),
    };

    console.log("LE BET", bet);
    try {
      const response = fetch("http://localhost:9000/bet", {
        method: "POST",
        headers: {
          "Content-Type": "application/json; charset=UTF-8",
        },
        mode: "no-cors",
        redirect: "manual",
        body: JSON.stringify(bet),
      });
      console.log("KJSDKFQLSKJMLKJ");
      const result = response.json();
      alert(result.note || result.error);
      console.log("result: ", result);
      return result;
    } catch (error) {
      console.error("Error registering bet:", error);
    }
  });

  getBlockchainButton.addEventListener("click", async () => {
    console.log("getblockcahin");
    try {
      const response = await fetch("http://localhost:9000/blockchain");
      const blockchain = await response.json();
      blockchainData.textContent = JSON.stringify(blockchain, null, 2);
    } catch (error) {
      console.error("Error fetching blockchain:", error);
    }
  });

  mineBlockButton.addEventListener("click", async () => {
    try {
      const response = await fetch("http://localhost:9000/mine");
      const result = await response.json();
      mineResult.textContent = JSON.stringify(result, null, 2);
    } catch (error) {
      console.error("Error mining block:", error);
    }
  });
});
