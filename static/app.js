// Check if the browser supports notifications
if ('Notification' in window) {
    // Request permission from the user
    Notification.requestPermission().then(function(permission) {
        if (permission === "granted") {
            console.log("Notification permission granted.");
            subscribeUserToPushNotifications();  // Call the function to subscribe
        } else {
            console.log("Notification permission denied.");
        }
    });
}

function subscribeUserToPushNotifications() {
    // Check if service worker and PushManager are available
    if ('serviceWorker' in navigator && 'PushManager' in window) {
        navigator.serviceWorker.ready.then(function(registration) {
            registration.pushManager.subscribe({
                userVisibleOnly: true,
                applicationServerKey: 'BDwYEyB3V2_NMzEgcGascHE3PUSQVPob7mnKyA5Qf8gzUqBWKDqlJQ_LujMSPbuYoWHH64pGKSnNJFtCANbTETM'  // Replace with your actual VAPID public key
            }).then(function(subscription) {
                console.log('User is subscribed:', subscription);

                // Send the subscription to your backend server to save it
                fetch('/subscribe', {
                    method: 'POST',
                    body: JSON.stringify(subscription),
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
            }).catch(function(error) {
                console.error('Error subscribing to push notifications:', error);
            });
        });
    }
}

// Register the service worker
if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/sw.js')
        .then(function(registration) {
            console.log('Service Worker registered with scope:', registration.scope);
        })
        .catch(function(error) {
            console.log('Service Worker registration failed:', error);
        });
}
