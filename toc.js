// Populate the sidebar
//
// This is a script, and not included directly in the page, to control the total size of the book.
// The TOC contains an entry for each page, so if each page includes a copy of the TOC,
// the total size of the page becomes O(n**2).
class MDBookSidebarScrollbox extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        this.innerHTML = '<ol class="chapter"><li class="chapter-item "><a href="carapace-spec.html"><strong aria-hidden="true">1.</strong> carapace-spec</a><a class="toggle"><div>❱</div></a></li><li><ol class="section"><li class="chapter-item "><a href="carapace-spec/usage.html"><strong aria-hidden="true">1.1.</strong> Usage</a></li><li class="chapter-item "><a href="carapace-spec/command.html"><strong aria-hidden="true">1.2.</strong> Command</a><a class="toggle"><div>❱</div></a></li><li><ol class="section"><li class="chapter-item "><a href="carapace-spec/command/name.html"><strong aria-hidden="true">1.2.1.</strong> Name</a></li><li class="chapter-item "><a href="carapace-spec/command/aliases.html"><strong aria-hidden="true">1.2.2.</strong> Aliases</a></li><li class="chapter-item "><a href="carapace-spec/command/description.html"><strong aria-hidden="true">1.2.3.</strong> Description</a></li><li class="chapter-item "><a href="carapace-spec/command/group.html"><strong aria-hidden="true">1.2.4.</strong> Group</a></li><li class="chapter-item "><a href="carapace-spec/command/hidden.html"><strong aria-hidden="true">1.2.5.</strong> Hidden</a></li><li class="chapter-item "><a href="carapace-spec/command/flags.html"><strong aria-hidden="true">1.2.6.</strong> Flags</a></li><li class="chapter-item "><a href="carapace-spec/command/persistentFlags.html"><strong aria-hidden="true">1.2.7.</strong> PersistentFlags</a></li><li class="chapter-item "><a href="carapace-spec/command/exclusiveFlags.html"><strong aria-hidden="true">1.2.8.</strong> ExclusiveFlags</a></li><li class="chapter-item "><a href="carapace-spec/command/completion.html"><strong aria-hidden="true">1.2.9.</strong> Completion</a><a class="toggle"><div>❱</div></a></li><li><ol class="section"><li class="chapter-item "><a href="carapace-spec/command/completion/flag.html"><strong aria-hidden="true">1.2.9.1.</strong> Flag</a></li><li class="chapter-item "><a href="carapace-spec/command/completion/positional.html"><strong aria-hidden="true">1.2.9.2.</strong> Positional</a></li><li class="chapter-item "><a href="carapace-spec/command/completion/positionalAny.html"><strong aria-hidden="true">1.2.9.3.</strong> PositionalAny</a></li><li class="chapter-item "><a href="carapace-spec/command/completion/dash.html"><strong aria-hidden="true">1.2.9.4.</strong> Dash</a></li><li class="chapter-item "><a href="carapace-spec/command/completion/dashAny.html"><strong aria-hidden="true">1.2.9.5.</strong> DashAny</a></li></ol></li><li class="chapter-item "><a href="carapace-spec/command/parsing.html"><strong aria-hidden="true">1.2.10.</strong> Parsing</a></li><li class="chapter-item "><a href="carapace-spec/command/run.html"><strong aria-hidden="true">1.2.11.</strong> Run</a></li><li class="chapter-item "><a href="carapace-spec/command/commands.html"><strong aria-hidden="true">1.2.12.</strong> Commands</a></li></ol></li><li class="chapter-item "><a href="carapace-spec/values.html"><strong aria-hidden="true">1.3.</strong> Values</a></li><li class="chapter-item "><a href="carapace-spec/macros.html"><strong aria-hidden="true">1.4.</strong> Macros</a><a class="toggle"><div>❱</div></a></li><li><ol class="section"><li class="chapter-item "><a href="carapace-spec/macros/core.html"><strong aria-hidden="true">1.4.1.</strong> Core</a></li><li class="chapter-item "><a href="carapace-spec/macros/modifier.html"><strong aria-hidden="true">1.4.2.</strong> Modifier</a></li><li class="chapter-item "><a href="carapace-spec/macros/custom.html"><strong aria-hidden="true">1.4.3.</strong> Custom</a></li></ol></li><li class="chapter-item "><a href="carapace-spec/variables.html"><strong aria-hidden="true">1.5.</strong> Variables</a></li></ol></li></ol>';
        // Set the current, active page, and reveal it if it's hidden
        let current_page = document.location.href.toString().split("#")[0].split("?")[0];
        if (current_page.endsWith("/")) {
            current_page += "index.html";
        }
        var links = Array.prototype.slice.call(this.querySelectorAll("a"));
        var l = links.length;
        for (var i = 0; i < l; ++i) {
            var link = links[i];
            var href = link.getAttribute("href");
            if (href && !href.startsWith("#") && !/^(?:[a-z+]+:)?\/\//.test(href)) {
                link.href = path_to_root + href;
            }
            // The "index" page is supposed to alias the first chapter in the book.
            if (link.href === current_page || (i === 0 && path_to_root === "" && current_page.endsWith("/index.html"))) {
                link.classList.add("active");
                var parent = link.parentElement;
                if (parent && parent.classList.contains("chapter-item")) {
                    parent.classList.add("expanded");
                }
                while (parent) {
                    if (parent.tagName === "LI" && parent.previousElementSibling) {
                        if (parent.previousElementSibling.classList.contains("chapter-item")) {
                            parent.previousElementSibling.classList.add("expanded");
                        }
                    }
                    parent = parent.parentElement;
                }
            }
        }
        // Track and set sidebar scroll position
        this.addEventListener('click', function(e) {
            if (e.target.tagName === 'A') {
                sessionStorage.setItem('sidebar-scroll', this.scrollTop);
            }
        }, { passive: true });
        var sidebarScrollTop = sessionStorage.getItem('sidebar-scroll');
        sessionStorage.removeItem('sidebar-scroll');
        if (sidebarScrollTop) {
            // preserve sidebar scroll position when navigating via links within sidebar
            this.scrollTop = sidebarScrollTop;
        } else {
            // scroll sidebar to current active section when navigating via "next/previous chapter" buttons
            var activeSection = document.querySelector('#sidebar .active');
            if (activeSection) {
                activeSection.scrollIntoView({ block: 'center' });
            }
        }
        // Toggle buttons
        var sidebarAnchorToggles = document.querySelectorAll('#sidebar a.toggle');
        function toggleSection(ev) {
            ev.currentTarget.parentElement.classList.toggle('expanded');
        }
        Array.from(sidebarAnchorToggles).forEach(function (el) {
            el.addEventListener('click', toggleSection);
        });
    }
}
window.customElements.define("mdbook-sidebar-scrollbox", MDBookSidebarScrollbox);
