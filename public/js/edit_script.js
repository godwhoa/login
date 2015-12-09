$(".input").focusin(function() {
    $(this).find("span").animate({
        "opacity": "0"
    }, 200);
});

$(".input").focusout(function() {
    $(this).find("span").animate({
        "opacity": "1"
    }, 300);
});



var correct = function() {
    $('.login').find(".submit i").removeAttr('class').addClass("fa fa-check").css({
        "color": "#fff"
    });
    $(".submit").css({
        "background": "#2ecc71",
        "border-color": "#2ecc71"
    });
    $(".feedback").show().animate({
        "opacity": "1",
        "bottom": "-80px"
    }, 400);
    $("input").css({
        "border-color": "#2ecc71"
    });
    window.location = "/profile"
}

var incorrect = function() {
    $('.login').find(".submit i").removeAttr('class').addClass("fa fa-times").css({
        "color": "#fff"
    });
    $(".submit").css({
        "background": "#cc2e3a",
        "border-color": "#cc2e3a"
    });
    $(".nfeedback").show().animate({
        "opacity": "1",
        "bottom": "-80px"
    }, 400)
    $("input").css({
        "border-color": "#cc2e3a"
    });
}


var post = function() {
    var formData = new FormData($(".login")[0]);


    $.ajax({
        url: window.location.pathname,
        data: formData,
        processData: false,
        contentType: false,
        type: 'POST',
        success: function(data) {
            console.log(data)
            if (data == "pos") {
                correct()

            } else {
                incorrect()
            }
        }
    });
}

$(".login").submit(function() {
    post()
    return false;
});
