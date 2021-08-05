// this is a helper to paste in the browser to pull out urls from aws side nav
function getAwsUrlsFromPage(className) {
  return $$(`.${className} a`)
    .map(element => {
      let url;
      try {
        url = new URL(element.href);
        url.searchParams.delete('region');
        url = url.pathname + url.search + url.hash;
      } catch {
        url = element.href;
      }
      return {
        id: element.innerText.replaceAll(' ', '').replaceAll('-', '').toLowerCase(),
        name: element.innerText,
        url,
      };
    })
    .filter(entry => !!entry.id);
}
