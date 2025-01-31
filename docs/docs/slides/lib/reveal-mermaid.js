window.setTimeout(function(){
    mermaids = document.getElementsByClassName('language-mermaid');
    let i=0;
    for (let el of mermaids) {
      newNode = document.createElement('div');
      newNode.innerHTML = el.innerHTML;
      newNode.classList.add('mermaid');
      newNode.id = 'mermaid-' + i;
      // This is temporarily placed to documentElement because styling issues.
      el.parentNode.setAttribute('data-id', newNode.id);
      document.documentElement.appendChild(newNode);
      i++;
    }
  
    // Generate mermaid diagrams
    try {
      mermaid.init({
        theme: 'dark',
        themeCSS: '.node rect { fill: red; }',
        flowchart:{
          useMaxWidth:false
        }
      }, '.mermaid');
    } catch (e) {}
  
    // Move new nodes in place of the old code
    while (mermaids.length > 0) {
      let el = mermaids.item(0);
      oldNode = el.parentNode;
      id = oldNode.getAttribute('data-id')
      console.log(id);
      oldNode.parentNode.replaceChild(document.getElementById(id), oldNode);
    }
    window.dispatchEvent(new Event('resize'));
  }, 5000);
  