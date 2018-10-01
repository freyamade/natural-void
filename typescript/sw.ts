self.addEventListener('install', (e: InstallEvent) => {
  e.waitUntil(
    caches.open('natural-void').then((cache: Cache) => {
      // Only initially cache the homepage. We're going to take a live-first approach to the service worker though
      return cache.addAll([
        '/',
        '/?source=pwa',
      ]);
    })
  );
});

// Handle fetching pages using online-first approach
self.addEventListener('fetch', (e: FetchEvent) => {
  console.log(e.request);
  // e.respondWith(
  //   fetch(e.request).then((response) => {
  //     if (response.ok && )
  //   });
  // );
});
