self.addEventListener('push', function(event) {
    const data = event.data ? event.data.text() : 'Default message';
    event.waitUntil(
      self.registration.showNotification('Push Notification', {
        body: data,
      })
    );
  });
  