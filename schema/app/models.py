from django.db import models


class User(models.Model):
    id = models.BigAutoField(primary_key=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    name = models.CharField(max_length=300)
    email = models.EmailField(null=False)
    password = models.CharField(max_length=300)
