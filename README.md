# BytePayments
BytePayments is a self hosted crypto payment gateway for accepting crypto funds directly to your wallet. For now it only supports TRX.

## Features :
1. Accept TRX Payment.
2. Create Payment.
3. Cancel any created payment.
4. List the available currencies (It's an array though but we have only trx for now, planned to add more in future).
5. Check a created payment status (completed,pending,cancelled).
6. Set the percentage of amount that is okay to be paid to mark the order as completed (eg : 95% payment marks the order as completed).
7. Handle Overpaid and Underpaid senario.
3. Send payment invoice directly to the users email after done.
4. After payment done move the funds to your main master wallet (Gas Fees Auto Calculated).

## Tech Stack :
1. Go (the goat).
2. Fiber (web framework based on fasthttp,net/http kinda slow)
3. Gorm (ORM library).
4. MySQL the OG DB.



### NOTE : This project has been tested internally and currently in alpha. Use at your own risk cause this can cause fund loss. Also if Crypto Currency maybe banned or not recognized in your country so use it at your own risk, we are not responsible for any loss or legal drama.
