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
    const options = {
        body: event.data ? event.data.text() : 'Here is a notification!',
        icon: '/icons/icon-192x192.png',
        badge: '/icons/icon-192x192.png',
    };

    event.waitUntil(
        self.registration.showNotification('Push Notification', options)
    );
});

self.addEventListener('notificationclick', (event) => {
    console.log('Notification clicked:', event);
    event.notification.close();
    event.waitUntil(
        clients.openWindow('/')
    );
});
