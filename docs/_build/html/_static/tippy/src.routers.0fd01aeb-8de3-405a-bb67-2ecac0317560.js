selector_to_html = {"a[href=\"#module-src.routers.classes\"]": "<h2 class=\"tippy-header\" style=\"margin-top: 0;\">src.routers.classes module<a class=\"headerlink\" href=\"#module-src.routers.classes\" title=\"Permalink to this heading\">#</a></h2>", "a[href=\"#src-routers-package\"]": "<h1 class=\"tippy-header\" style=\"margin-top: 0;\">src.routers package<a class=\"headerlink\" href=\"#src-routers-package\" title=\"Permalink to this heading\">#</a></h1><h2>Submodules<a class=\"headerlink\" href=\"#submodules\" title=\"Permalink to this heading\">#</a></h2>", "a[href=\"#submodules\"]": "<h2 class=\"tippy-header\" style=\"margin-top: 0;\">Submodules<a class=\"headerlink\" href=\"#submodules\" title=\"Permalink to this heading\">#</a></h2>", "a[href=\"#module-src.routers\"]": "<h2 class=\"tippy-header\" style=\"margin-top: 0;\">Module contents<a class=\"headerlink\" href=\"#module-src.routers\" title=\"Permalink to this heading\">#</a></h2>", "a[href=\"#src.routers.classes.all_classes\"]": "<dt class=\"sig sig-object py\" id=\"src.routers.classes.all_classes\">\n<em class=\"property\"><span class=\"k\"><span class=\"pre\">async</span></span><span class=\"w\"> </span></em><span class=\"sig-prename descclassname\"><span class=\"pre\">src.routers.classes.</span></span><span class=\"sig-name descname\"><span class=\"pre\">all_classes</span></span><span class=\"sig-paren\">(</span><em class=\"sig-param\"><span class=\"n\"><span class=\"pre\">page</span></span><span class=\"p\"><span class=\"pre\">:</span></span><span class=\"w\"> </span><span class=\"n\"><span class=\"pre\">int</span></span><span class=\"w\"> </span><span class=\"o\"><span class=\"pre\">=</span></span><span class=\"w\"> </span><span class=\"default_value\"><span class=\"pre\">1</span></span></em>, <em class=\"sig-param\"><span class=\"n\"><span class=\"pre\">page_size</span></span><span class=\"p\"><span class=\"pre\">:</span></span><span class=\"w\"> </span><span class=\"n\"><span class=\"pre\">int</span></span><span class=\"w\"> </span><span class=\"o\"><span class=\"pre\">=</span></span><span class=\"w\"> </span><span class=\"default_value\"><span class=\"pre\">10</span></span></em>, <em class=\"sig-param\"><span class=\"n\"><span class=\"pre\">db</span></span><span class=\"p\"><span class=\"pre\">:</span></span><span class=\"w\"> </span><span class=\"n\"><span class=\"pre\">AsyncSession</span></span><span class=\"w\"> </span><span class=\"o\"><span class=\"pre\">=</span></span><span class=\"w\"> </span><span class=\"default_value\"><span class=\"pre\">Depends(get_db)</span></span></em><span class=\"sig-paren\">)</span></dt><dd></dd>"}
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
