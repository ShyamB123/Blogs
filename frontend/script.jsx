// Remove showForm function, as modals are now used

function toggleBlog(id) {
  const el = document.getElementById(id);
  if (el) {
    if (el.style.display === 'none' || el.style.display === '') {
      el.style.display = 'block';
      el.classList.add('blog-details-animate');
      setTimeout(() => el.classList.remove('blog-details-animate'), 400);
    } else {
      el.style.display = 'none';
    }
  }
}

// Add a little animation CSS via JS if not present
(function addAnimStyles() {
  if (!document.getElementById('extra-anim-style')) {
    const style = document.createElement('style');
    style.id = 'extra-anim-style';
    style.innerHTML = `
      .form-highlight {
        box-shadow: 0 0 0 3px #2d6cdf44;
        transition: box-shadow 0.4s;
      }
      .blog-details-animate {
        animation: blogExpand 0.4s;
      }
      @keyframes blogExpand {
        from { opacity: 0; transform: scaleY(0.95); }
        to { opacity: 1; transform: scaleY(1); }
      }
    `;
    document.head.appendChild(style);
  }
})();
