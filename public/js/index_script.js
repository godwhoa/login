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
    var user = $("#user").val()
    var pass = $("#pass").val()
    console.log(user+" "+pass)
    $.ajax({
            type: 'POST',
            url: '/',
            data: {
                user: user,
                pass: pass
            },
            success: function(data) {
                console.log(data)
                if (data == "neg") {
                    incorrect()
                }else{
                    correct()
                }
            }
    });
}


$(".login").submit(function() {
    post()
    return false;
});

$(".reg").click(function() {
    window.location = "/register"
    return false;
});
