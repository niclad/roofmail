console.info("Welcome to Roofmail!");

let likeBtn = document.getElementById("like-btn");
let dislikeBtn = document.getElementById("dislike-btn");
let likeIcon = document.getElementById("like-icon");
let dislikeIcon = document.getElementById("dislike-icon");

function likeMouseEnter() {
    likeBtn.classList.remove("btn-primary");
    likeBtn.classList.add("btn-success");
    likeIcon.classList.remove("bi-hand-thumbs-up");
    likeIcon.classList.add("bi-hand-thumbs-up-fill");
}

function likeMouseLeave() {
    likeBtn.classList.remove("btn-success");
    likeBtn.classList.add("btn-primary");
    likeIcon.classList.remove("bi-hand-thumbs-up-fill");
    likeIcon.classList.add("bi-hand-thumbs-up");
}

function dislikeMouseEnter() {
    dislikeBtn.classList.remove("btn-primary");
    dislikeBtn.classList.add("btn-danger");
    dislikeIcon.classList.remove("bi-hand-thumbs-down");
    dislikeIcon.classList.add("bi-hand-thumbs-down-fill");
}

function dislikeMouseLeave() {
    dislikeBtn.classList.remove("btn-danger");
    dislikeBtn.classList.add("btn-primary");
    dislikeIcon.classList.remove("bi-hand-thumbs-down-fill");
    dislikeIcon.classList.add("bi-hand-thumbs-down");
}

function toggleLikeButton() {
    likeBtn.classList.remove("btn-success");
    likeBtn.classList.add("btn-primary");

    likeIcon.classList.remove("bi-hand-thumbs-up-fill");
    likeIcon.classList.add("bi-hand-thumbs-up");

    likeBtn.addEventListener("mouseenter", likeMouseEnter);
    likeBtn.addEventListener("mouseleave", likeMouseLeave);
}

toggleLikeButton();

function toggleDislikeButton() {
    dislikeBtn.classList.remove("btn-danger");
    dislikeBtn.classList.add("btn-primary");

    dislikeIcon.classList.remove("bi-hand-thumbs-down-fill");
    dislikeIcon.classList.add("bi-hand-thumbs-down");

    dislikeBtn.addEventListener("mouseenter", dislikeMouseEnter);
    dislikeBtn.addEventListener("mouseleave", dislikeMouseLeave);
}

toggleDislikeButton();

let liked = null;

function submitLike(isLiked) {
    if (liked === isLiked) {
        return; // No change in like status, do nothing
    }

    liked = isLiked;


    fetch("/like", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ Liked: isLiked })
    }).then(response => {
        if (response.ok) {
            // console.debug("Like status submitted successfully.");
        } else {
            console.error("Failed to submit like status.");
        }
    }).catch(error => {
        console.error("Error submitting like status:", error);
    });
}

likeBtn.addEventListener("click", () => {
    submitLike(true);

    likeBtn.removeEventListener("mouseenter", likeMouseEnter);
    likeBtn.removeEventListener("mouseleave", likeMouseLeave);

    toggleDislikeButton();

    likeBtn.classList.remove("btn-primary");
    likeBtn.classList.add("btn-success");
    likeIcon.classList.remove("bi-hand-thumbs-up");
    likeIcon.classList.add("bi-hand-thumbs-up-fill");
});

dislikeBtn.addEventListener("click", () => {
    submitLike(false);

    dislikeBtn.removeEventListener("mouseenter", dislikeMouseEnter);
    dislikeBtn.removeEventListener("mouseleave", dislikeMouseLeave);

    toggleLikeButton();

    dislikeBtn.classList.remove("btn-primary");
    dislikeBtn.classList.add("btn-danger");
    dislikeIcon.classList.remove("bi-hand-thumbs-down");
    dislikeIcon.classList.add("bi-hand-thumbs-down-fill");
});
