(() => {
    const link = document.getElementById('new-post');
    if (!link) return;

    const updateLink = () => {
        const now = new Date();
        const date = [
            now.getFullYear(),
            String(now.getMonth() + 1).padStart(2, '0'),
            String(now.getDate()).padStart(2, '0'),
        ].join('-');
        const url = new URL(link.href);
        url.searchParams.set('filename', `${date}-title.md`);
        url.searchParams.set('value', `---\nlayout: post\ntitle: \ndate: ${date}\n---\n`);
        link.href = url.href;
    };

    updateLink();
    link.addEventListener('focus', updateLink);
    link.addEventListener('pointerdown', updateLink);
    link.addEventListener('click', updateLink);
})();
