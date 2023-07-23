selector_to_html = {"a[href=\"modules.html\"]": "<h1 class=\"tippy-header\" style=\"margin-top: 0;\">Documentation<a class=\"headerlink\" href=\"#documentation\" title=\"Permalink to this heading\">#</a></h1>", "a[href=\"src.html\"]": "<h1 class=\"tippy-header\" style=\"margin-top: 0;\">src package<a class=\"headerlink\" href=\"#src-package\" title=\"Permalink to this heading\">#</a></h1><h2>Subpackages<a class=\"headerlink\" href=\"#subpackages\" title=\"Permalink to this heading\">#</a></h2>", "a[href=\"#flamethrower-dnd-3-5e-rest-api\"]": "<h1 class=\"tippy-header\" style=\"margin-top: 0;\">Flamethrower - DnD &amp; 3.5e REST API<a class=\"headerlink\" href=\"#flamethrower-dnd-3-5e-rest-api\" title=\"Permalink to this heading\">#</a></h1>"}
skip_classes = ["headerlink", "sd-stretched-link"]

window.onload = function () {
    for (const [select, tip_html] of Object.entries(selector_to_html)) {
        const links = document.querySelectorAll(` ${select}`);
        for (const link of links) {
            if (skip_classes.some(c => link.classList.contains(c))) {
                continue;
            }

            tippy(link, {
                content: tip_html,
                allowHTML: true,
                arrow: true,
                placement: 'auto-start', maxWidth: 500, interactive: false,

            });
        };
    };
    console.log("tippy tips loaded!");
};
