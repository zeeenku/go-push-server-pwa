self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open('pwa-cache').then((cache) => {
            return cache.addAll([
                '/',
                '/index.html',
                '/manifest.json',
                '/icons/icon-192x192.png',
                '/icons/icon-512x512.png'
            ]);
        })
    );
    console.log('Service Worker installing...');
});

self.addEventListener('push', (event) => {
    console.log('Push notification received:', event);

    const notif = event.data ? event.data.text() : 'Here is a notification!';
    const options = {
        body: notif,
        icon: '/icons/icon-192x192.png', // Ensure this path points to a valid icon file
        badge: '/icons/icon-192x192.png', // Badge is also important for iOS
    };

    // Ensure the notification is shown even if no body data exists
    event.waitUntil(
        self.registration.showNotification("You might need this quote", options)
    );
});


self.addEventListener('notificationclick', (event) => {
    console.log('Notification clicked:', event);
    event.notification.close();
    event.waitUntil(
        clients.openWindow('/')
    );
});
