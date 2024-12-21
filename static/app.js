// Check if the browser supports notifications
if ('Notification' in window) {
    // Request permission to show notifications
    Notification.requestPermission().then(function(permission) {
        if (permission === "granted") {
            console.log("Notification permission granted.");
            subscribeUserToPushNotifications();  // Proceed with subscription
        } else {
            console.log("Notification permission denied.");
        }
    });
} else {
    console.log("Notifications are not supported by this browser.");
}

// Function to subscribe the user to push notifications
function subscribeUserToPushNotifications() {
    if ('serviceWorker' in navigator && 'PushManager' in window) {
        navigator.serviceWorker.register('/sw.js')
            .then(function(registration) {
                console.log('Service Worker registered with scope:', registration.scope);
                return registration.pushManager.subscribe({
                    userVisibleOnly: true,  // Ensure the user can see the notification
                    applicationServerKey: urlBase64ToUint8Array('BDwYEyB3V2_NMzEgcGascHE3PUSQVPob7mnKyA5Qf8gzUqBWKDqlJQ_LujMSPbuYoWHH64pGKSnNJFtCANbTETM')  // Replace with your VAPID public key
                });
            })
            .then(function(subscription) {
                console.log('User is subscribed:', subscription);

                // Send the subscription object to the server to store it
                fetch('/subscribe', {
                    method: 'POST',
                    body: JSON.stringify(subscription),
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
            })
            .catch(function(error) {
                console.error('Error subscribing to push notifications:', error);
            });
    } else {
        console.log("Service Worker or PushManager not supported.");
    }
}

// Function to convert VAPID public key to Uint8Array (base64url to binary)
function urlBase64ToUint8Array(base64String) {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);  // Add padding to make it 4-byte aligned
    const base64 = (base64String + padding)
        .replace(/\-/g, '+')   // Replace URL-safe characters
        .replace(/_/g, '/');    // Replace URL-safe characters
    const rawData = atob(base64);  // Decode base64 to raw data
    const outputArray = new Uint8Array(rawData.length);  // Create a Uint8Array of the same length

    for (let i = 0; i < rawData.length; ++i) {
        outputArray[i] = rawData.charCodeAt(i);  // Populate the array with the raw data
    }

    return outputArray;  // Return the binary data
}

